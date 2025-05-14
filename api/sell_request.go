package api

import (
	"errors"
	"fmt"
	"net/http"
	db "p2platform/db/sqlc"
	appErr "p2platform/errors"
	"p2platform/notification/kafka"
	"p2platform/notification/model"
	util "p2platform/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type createSellRequestQuery struct {
	SpaceId string `form:"space_id" binding:"required,uuid"`
}
type createSellRequest struct {
	SellTotalAmount  int64  `json:"sell_total_amount" binding:"required,min=1"`
	SellMoneySource  string `json:"sell_money_source" binding:"required,source"`
	CurrencyFrom     string `json:"currency_from" binding:"required,currency"`
	CurrencyTo       string `json:"currency_to" binding:"required,currency"`
	SellAmountByCard int64  `json:"sell_amount_by_card" binding:"gte=0"`
	SellAmountByCash int64  `json:"sell_amount_by_cash" binding:"gte=0"`
	SellExchangeRate int64  `json:"sell_exchange_rate" binding:"required,gt=0"`
	Comment          string `json:"comment" binding:"omitempty"`
}

type sellRequestResponse struct {
	SellReqID        int32     `json:"sell_req_id"`
	SellTotalAmount  int64     `json:"sell_total_amount"`
	SellMoneySource  string    `json:"sell_money_source"`
	CurrencyFrom     string    `json:"currency_from"`
	CurrencyTo       string    `json:"currency_to"`
	TgUsername       string    `json:"tg_username"`
	SellByCard       bool      `json:"sell_by_card"`
	SellAmountByCard int64     `json:"sell_amount_by_card"`
	SellByCash       bool      `json:"sell_by_cash"`
	SellAmountByCash int64     `json:"sell_amount_by_cash"`
	SellExchangeRate int64     `json:"sell_exchange_rate"`
	Comment          string    `json:"comment"`
	IsActual         bool      `json:"is_actual"`
	CreatedAt        time.Time `json:"created_at"`
}

type ErrResponse struct {
	Error string `json:"error" example:"something went wrong"`
}

// createSellRequest godoc
//
//	@Summary      Create a new sell request
//	@Description  Create a new sell request with Telegram ID and username
//	@Tags         sell-request
//	@Accept       json
//	@Produce      json
//	@Param        request  body      createSellRequest  true  "Create sell request"
//	@Success      200      {object}  sellRequestResponse
//	@Failure      400      {object}  ErrResponse
//	@Failure      404      {object}  ErrResponse
//	@Failure      500      {object}  ErrResponse
//	@Router       /sell-request [post]
func (server *Server) createSellRequest(ctx *gin.Context) {
	var uri createSellRequestQuery
	var req createSellRequest
	sellByCard := false
	sellByCash := false
	if err := ctx.ShouldBindQuery(&uri); err != nil {
		ctx.JSON(appErr.ErrInvalidUri.Status, ErrorResponse(appErr.ErrInvalidUri))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(appErr.ErrInvalidPayload.Status, ErrorResponse(appErr.ErrInvalidPayload))
		return
	}
	telegramId, ok := GetTelegramIDFromContext(ctx)
	if !ok {
		return
	}
	uid, err := uuid.Parse(uri.SpaceId)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidUUID.Status, ErrorResponse(appErr.ErrInvalidUUID))
		return
	}
	argIsUser := db.IsUserInSpaceParams{
		UserID:  telegramId,
		SpaceID: uid,
	}
	isMember, err := server.store.IsUserInSpace(ctx, argIsUser)
	if err != nil || !isMember {
		ctx.JSON(appErr.ErrForbidden.Status, ErrorResponse(appErr.ErrForbidden))
		return
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
	if !(req.SellAmountByCard+req.SellAmountByCash == req.SellTotalAmount) {
		ctx.JSON(appErr.ErrSellAmountMismatch.Status, ErrorResponse(appErr.ErrSellAmountMismatch))
		return
	}
	if req.SellAmountByCard > 0 {
		sellByCard = true
	}
	if req.SellAmountByCash > 0 {
		sellByCash = true
	}

	arg := db.CreateSellRequestParams{
		SpaceID:          uid,
		SellTotalAmount:  req.SellTotalAmount,
		SellMoneySource:  req.SellMoneySource,
		CurrencyFrom:     req.CurrencyFrom,
		CurrencyTo:       req.CurrencyTo,
		TelegramID:       telegramId,
		TgUsername:       user.TgUsername,
		SellByCard:       util.ToPgBool(sellByCard),
		SellAmountByCard: util.ToPgInt(req.SellAmountByCard),
		SellByCash:       util.ToPgBool(sellByCash),
		SellAmountByCash: util.ToPgInt(req.SellAmountByCash),
		SellExchangeRate: util.ToPgInt(req.SellExchangeRate),
		Comment:          req.Comment,
	}
	sellRequest, err := server.store.CreateSellRequest(ctx, arg)
	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
		return
	}
	result := sellRequestResponse{
		SellReqID:        sellRequest.SellReqID,
		SellTotalAmount:  sellRequest.SellTotalAmount,
		SellMoneySource:  sellRequest.SellMoneySource,
		CurrencyFrom:     sellRequest.CurrencyFrom,
		CurrencyTo:       sellRequest.CurrencyTo,
		TgUsername:       sellRequest.TgUsername,
		SellByCard:       sellRequest.SellByCard.Bool,
		SellAmountByCard: sellRequest.SellAmountByCard.Int64,
		SellByCash:       sellRequest.SellByCash.Bool,
		SellAmountByCash: sellRequest.SellAmountByCash.Int64,
		SellExchangeRate: sellRequest.SellExchangeRate.Int64,
		Comment:          sellRequest.CurrencyTo,
		IsActual:         sellRequest.IsActual.Bool,
		CreatedAt:        sellRequest.CreatedAt,
	}
	ctx.JSON(http.StatusOK, result)

	_ = kafka.Publish(server.producer, "notifications", model.NotifictationMessage{
		TelegramId: telegramId,
		Message:    fmt.Sprintf("Sell request #%d успешно создан!", sellRequest.SellReqID),
		EventType:  "sell_request_created",
	})

}

type getSellRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getSellRequest(ctx *gin.Context) {
	var req getSellRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
	telegramId, exists := GetTelegramIDFromContext(ctx)
	if !exists {
		return
	}

	result, err := server.store.GetSellRequestTx(ctx, req.ID, telegramId)
	if err != nil {
		HandleAppError(ctx, err)
		return
	}

	if result.SellRequest.IsDeleted.Bool {
		ctx.JSON(appErr.ErrSellRequestDeleted.Status, ErrorResponse(appErr.ErrSellRequestDeleted))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

type listSellRequestBySpace struct {
	PageSize int32  `form:"page_size" binding:"required,min=5,max=10"`
	PageId   int32  `form:"page_id" binding:"required,min=1"`
	SpaceId  string `form:"space_id" binding:"required,uuid"`
}

func (server *Server) listSellRequests(ctx *gin.Context) {
	var req listSellRequestBySpace
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidQuery.Status, ErrorResponse(appErr.ErrInvalidQuery))
		return
	}
	telegramId, exists := GetTelegramIDFromContext(ctx)
	if !exists {
		return
	}

	uid, err := uuid.Parse(req.SpaceId)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidUUID.Status, ErrorResponse(appErr.ErrInvalidUUID))
		return
	}

	arg := db.ListSellRequeststTxParams{
		Limit:   req.PageSize,
		Offset:  (req.PageId - 1) * req.PageSize,
		SpaceId: uid,
	}
	sellRequests, err := server.store.ListSellRequeststTx(ctx, arg, telegramId)
	if err != nil {
		HandleAppError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, sellRequests)
}

func (server *Server) listMySellRequests(ctx *gin.Context) {
	var req listSellRequestBySpace
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidQuery.Status, ErrorResponse(appErr.ErrInvalidQuery))
		return
	}
	uid, err := uuid.Parse(req.SpaceId)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidUUID.Status, ErrorResponse(appErr.ErrInvalidUUID))
		return
	}
	telegramId, ok := GetTelegramIDFromContext(ctx)
	if !ok {
		return
	}
	arg := db.ListMySellRequestsTxParams{
		Limit:      req.PageSize,
		Offset:     (req.PageId - 1) * req.PageSize,
		TelegramId: telegramId,
		SpaceId:    uid,
	}
	sellRequests, err := server.store.ListMySellRequeststTx(ctx, arg)
	if err != nil {
		HandleAppError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, sellRequests)
}

type updateSellRequestUri struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type updateSellRequestJson struct {
	SellTotalAmount  *int64  `json:"sell_total_amount" binding:"omitempty,min=1"`
	SellMoneySource  *string `json:"sell_money_source" binding:"omitempty,source"`
	CurrencyFrom     *string `json:"currency_from" binding:"omitempty,currency"`
	CurrencyTo       *string `json:"currency_to" binding:"omitempty,currency"`
	SellAmountByCard *int64  `json:"sell_amount_by_card" binding:"omitempty,gte=0"`
	SellAmountByCash *int64  `json:"sell_amount_by_cash" binding:"omitempty,gte=0"`
	SellExchangeRate *int64  `json:"sell_exchange_rate" binding:"omitempty,min=1"`
	Comment          *string `json:"comment" binding:"omitempty,max=100"`
}

func (server *Server) updateSellRequest(ctx *gin.Context) {
	var reqUri updateSellRequestUri
	var reqJson updateSellRequestJson
	var totalAmountToCheckSum int64
	var amountByCardToCheckSum int64
	var amountByCashToCheckSum int64
	var sellByCard pgtype.Bool
	var sellByCash pgtype.Bool

	err := ctx.ShouldBindUri(&reqUri)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidUri.Status, ErrorResponse(appErr.ErrInvalidUri))
		return
	}
	err = ctx.ShouldBindJSON(&reqJson)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidPayload.Status, ErrorResponse(appErr.ErrInvalidPayload))
		return
	}

	telegramId, ok := GetTelegramIDFromContext(ctx)
	if !ok {
		return
	}

	buyRequests, err := server.store.ListBuyRequests(ctx, db.ListBuyRequestsParams{
		SellReqID: reqUri.ID,
		Limit:     1,
		Offset:    0,
	})
	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
		return
	}
	if len(buyRequests) > 0 {
		ctx.JSON(appErr.ErrSellRequestHasBuyRequests.Status, ErrorResponse(appErr.ErrSellRequestHasBuyRequests))
		return
	}

	sellRequest, err := server.store.GetSellRequestById(ctx, reqUri.ID)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsFound) {
			ctx.JSON(appErr.ErrSellRequestNotFound.Status, ErrorResponse(appErr.ErrSellRequestNotFound))
			return
		}
		ctx.JSON(appErr.ErrFailedToGetSellRequests.Status, ErrorResponse(appErr.ErrFailedToGetSellRequests))
		return
	}

	if telegramId != sellRequest.TelegramID {
		ctx.JSON(appErr.ErrNotSellRequestOwner.Status, ErrorResponse(appErr.ErrNotSellRequestOwner))
		return
	}

	if sellRequest.IsDeleted.Bool {
		ctx.JSON(appErr.ErrSellRequestDeleted.Status, ErrorResponse(appErr.ErrSellRequestDeleted))
		return
	}

	if reqJson.SellTotalAmount != nil {
		totalAmountToCheckSum = util.DerefInt64(reqJson.SellTotalAmount)
	} else {
		totalAmountToCheckSum = sellRequest.SellTotalAmount
	}
	if reqJson.SellAmountByCard != nil {
		amountByCardToCheckSum = util.DerefInt64(reqJson.SellAmountByCard)
	} else {
		amountByCardToCheckSum = sellRequest.SellAmountByCard.Int64
	}
	if reqJson.SellAmountByCash != nil {
		amountByCashToCheckSum = util.DerefInt64(reqJson.SellAmountByCash)
	} else {
		amountByCashToCheckSum = sellRequest.SellAmountByCash.Int64
	}

	if totalAmountToCheckSum != amountByCardToCheckSum+amountByCashToCheckSum {
		ctx.JSON(appErr.ErrSellAmountMismatch.Status, ErrorResponse(appErr.ErrSellAmountMismatch))
		return
	}

	if reqJson.SellAmountByCard != nil {
		if util.DerefInt64(reqJson.SellAmountByCard) > 0 {
			sellByCard = util.ToPgBool(true)
		} else {
			sellByCard = util.ToPgBool(false)
		}
	} else {
		sellByCard = sellRequest.SellByCard
	}

	if reqJson.SellAmountByCash != nil {
		if util.DerefInt64(reqJson.SellAmountByCash) > 0 {
			sellByCash = util.ToPgBool(true)
		} else {
			sellByCash = util.ToPgBool(false)
		}
	} else {
		sellByCash = sellRequest.SellByCash
	}

	arg := db.UpdateSellRequestParams{
		SellReqID: reqUri.ID,

		SellTotalAmount: pgtype.Int8{
			Int64: util.DerefInt64(reqJson.SellTotalAmount),
			Valid: reqJson.SellTotalAmount != nil,
		},
		SellMoneySource: pgtype.Text{
			String: util.DerefStr(reqJson.SellMoneySource),
			Valid:  reqJson.SellMoneySource != nil,
		},
		CurrencyFrom: pgtype.Text{
			String: util.DerefStr(reqJson.CurrencyFrom),
			Valid:  reqJson.CurrencyFrom != nil,
		},
		CurrencyTo: pgtype.Text{
			String: util.DerefStr(reqJson.CurrencyTo),
			Valid:  reqJson.CurrencyTo != nil,
		},

		SellByCard: sellByCard,

		SellAmountByCard: pgtype.Int8{
			Int64: util.DerefInt64(reqJson.SellAmountByCard),
			Valid: reqJson.SellAmountByCard != nil,
		},
		SellAmountByCash: pgtype.Int8{
			Int64: util.DerefInt64(reqJson.SellAmountByCash),
			Valid: reqJson.SellAmountByCash != nil,
		},

		SellByCash: sellByCash,

		SellExchangeRate: pgtype.Int8{
			Int64: util.DerefInt64(reqJson.SellExchangeRate),
			Valid: reqJson.SellExchangeRate != nil,
		},
		Comment: pgtype.Text{
			String: util.DerefStr(reqJson.Comment),
			Valid:  reqJson.Comment != nil,
		},
	}
	sellRequest, err = server.store.UpdateSellRequest(ctx, arg)
	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
		return
	}
	ctx.JSON(http.StatusOK, sellRequest)
}

type deleteSellRequestUri struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}
type deleteSellRequestResponse struct {
	IsDeleted bool `json:"is_deleted"`
}

func (server *Server) deleteSellRequest(ctx *gin.Context) {
	var req deleteSellRequestUri
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidUri.Status, ErrorResponse(appErr.ErrInvalidUri))
		return
	}

	telegramId, ok := GetTelegramIDFromContext(ctx)
	if !ok {
		return
	}

	sellRequest, err := server.store.GetSellRequestById(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsFound) {
			ctx.JSON(appErr.ErrSellRequestNotFound.Status, ErrorResponse(appErr.ErrSellRequestNotFound))
			return
		}
		ctx.JSON(appErr.ErrFailedToGetSellRequests.Status, ErrorResponse(appErr.ErrFailedToGetSellRequests))
		return
	}

	if telegramId != sellRequest.TelegramID {
		ctx.JSON(appErr.ErrNotSellRequestOwner.Status, ErrorResponse(appErr.ErrNotSellRequestOwner))
		return
	}

	isDeleted, err := server.store.DeleteSellRequestTx(ctx, req.ID)
	if err != nil {
		HandleAppError(ctx, err)
		return
	}
	result := deleteSellRequestResponse{
		IsDeleted: isDeleted,
	}
	ctx.JSON(http.StatusOK, result)
}
