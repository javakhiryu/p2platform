package api

import (
	"net/http"
	db "p2platform/db/sqlc"
	appErr "p2platform/errors"
	"p2platform/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createSpaceRequest struct {
	SpaceName   string `json:"name" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	Description string `json:"description"`
}

func (server *Server) createSpace(ctx *gin.Context) {
	var req createSpaceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(appErr.ErrInvalidPayload.Status, ErrorResponse(appErr.ErrInvalidPayload))
		return
	}
	telegramId, ok := GetTelegramIDFromContext(ctx)
	if !ok {
		return
	}
	uuid, err := uuid.NewRandom()
	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
	}
	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
	}

	result, err := server.store.CreateSpaceTx(ctx, db.CreateSpaceTxParams{
		SpaceID:        uuid,
		SpaceName:      req.SpaceName,
		HashedPassword: hashedPassword,
		Description:    req.Description,
		CreatorID:      telegramId,
	})
	if err != nil {
		HandleAppError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

type listSpacesRequest struct {
	LastSpaceName string `form:"last_space_name"`
	LastSpaceID   string `form:"last_space_id"`
	Limit         int    `form:"limit"`
}

type spaceResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	CreatorID   int64     `json:"creator_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type NextCursor struct {
	LastSpaceName string    `form:"last_space_name"`
	LastSpaceID   uuid.UUID `form:"last_space_id"`
}

type listSpacesResponse struct {
	Spaces     []spaceResponse `json:"spaces"`
	HasMore    bool            `json:"has_more"`
	NextCursor NextCursor      `json:"next_cursor,omitempty"`
}

func (server *Server) listSpaces(ctx *gin.Context) {
	var req listSpacesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(appErr.ErrInvalidQuery.Status, ErrorResponse(appErr.ErrInvalidQuery))
		return
	}

	// Установка лимита
	limit := 10
	if req.Limit > 0 && req.Limit <= 100 {
		limit = req.Limit
	}

	var spaces []db.Space
	var err error

	if req.LastSpaceName != "" && req.LastSpaceID != "" {
		// Запрос с курсором
		_, err := uuid.Parse(req.LastSpaceID)
		if err != nil {
			ctx.JSON(appErr.ErrInvalidUUID.Status, ErrorResponse(appErr.ErrInvalidUUID))
			return
		}

		arg := db.ListSpacesAfterCursorByNameAscParams{
			Limit:     int32(limit + 1),
			SpaceName: req.LastSpaceName,
			SpaceID:   req.LastSpaceID,
		}
		spaces, err = server.store.ListSpacesAfterCursorByNameAsc(ctx, arg)

	} else {
		// Первая загрузка без курсора
		spaces, err = server.store.ListFirstSpacesByNameAsc(ctx, int32(limit+1))
	}

	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
		return
	}

	hasMore := len(spaces) > limit
	if hasMore {
		spaces = spaces[:limit]
	}

	response := listSpacesResponse{
		Spaces:  make([]spaceResponse, len(spaces)),
		HasMore: hasMore,
	}

	for i, space := range spaces {
		response.Spaces[i] = spaceResponse{
			ID:          space.SpaceID,
			Name:        space.SpaceName,
			CreatorID:   space.CreatorID.Int64,
			Description: space.Description,
			CreatedAt:   space.CreatedAt,
			UpdatedAt:   space.UpdatedAt,
		}
	}

	if hasMore {
		lastSpace := spaces[len(spaces)-1]
		response.NextCursor = NextCursor{
			LastSpaceName: lastSpace.SpaceName,
			LastSpaceID:   lastSpace.SpaceID,
		}
	}

	ctx.JSON(http.StatusOK, response)
}

type GetSpaceRequest struct {
	SpaceID string `uri:"id" binding:"required,uuid"`
}

func (server *Server) getSpace(ctx *gin.Context) {
	var req GetSpaceRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(appErr.ErrInvalidUri.Status, ErrorResponse(appErr.ErrInvalidUri))
		return
	}
	uid, err := uuid.Parse(req.SpaceID)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidUUID.Status, ErrorResponse(appErr.ErrInvalidUUID))
		return
	}
	telegramID, ok := GetTelegramIDFromContext(ctx)
	if !ok {
		return
	}
	response, err := server.store.GetSpaceTx(ctx, uid, telegramID)
	if err != nil {
		HandleAppError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
