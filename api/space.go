package api

import (
	"errors"
	"fmt"
	"net/http"
	db "p2platform/db/sqlc"
	appErr "p2platform/errors"
	"p2platform/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createSpaceRequest struct {
	SpaceName   string `json:"name" binding:"required,alphanum"`
	Password    string `json:"password" binding:"required,min=6"`
	Description string `json:"description"`
}
type createSpaceResponse struct {
	SpaceID     uuid.UUID `json:"space_id"`
	SpaceName   string    `json:"name"`
	User        db.User   `json:"user"`
	Description string    `json:"description"`
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
	user, err := server.store.GetUser(ctx, telegramId)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsFound) {
			ctx.JSON(appErr.ErrUserNotFound.Status, ErrorResponse(appErr.ErrUserNotFound))
			return
		}
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
		return
	}
	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
	}
	arg := db.CreateSpaceParams{
		SpaceID:        uuid,
		SpaceName:      req.SpaceName,
		HashedPassword: hashedPassword,
		Description:    req.Description,
		CreatorID:      util.ToPgInt(telegramId),
	}
	space, err := server.store.CreateSpace(ctx, arg)
	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
	}
	result := createSpaceResponse{
		SpaceID:     space.SpaceID,
		SpaceName:   space.SpaceName,
		User:        user,
		Description: space.Description,
	}
	ctx.JSON(http.StatusOK, result)
}

type listSpacesRequest struct {
	LastSpaceName string `form:"last_space_name"`
	LastSpaceID   string `form:"last_space_id"`
	Limit         int    `form:"limit"`
}

type spaceResponse struct {
	ID          uuid.UUID    `json:"id"`
	Name        string    `json:"name"`
	CreatorID   int64    `json:"creator_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type listSpacesResponse struct {
	Spaces     []spaceResponse `json:"spaces"`
	HasMore    bool            `json:"has_more"`
	NextCursor string          `json:"next_cursor,omitempty"`
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
        spaces, err = server.store.ListFirstSpacesByNameAsc(ctx, int32(limit + 1))
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
        response.NextCursor = fmt.Sprintf("%s,%s", lastSpace.SpaceName, lastSpace.SpaceID.String())
    }

    ctx.JSON(http.StatusOK, response)
}
