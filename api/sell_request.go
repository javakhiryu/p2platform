package api

import (
	"net/http"
	db "p2platform/db/sqlc"
	util "p2platform/util"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type createSellRequest struct {
	SellTotalAmount  int64  `json:"sell_total_amount" binding:"required,min=1"`
	CurrencyFrom     string `json:"currency_from" binding:"required,currency"`
	CurrencyTo       string `json:"currency_to" binding:"required,currency"`
	TgUsername       string `json:"tg_username" binding:"required"`
	SellByCard       bool   `json:"sell_by_card" binding:"required"`
	SellAmountByCard int64  `json:"sell_amount_by_card" binding:"required,gte=0"`
	SellByCash       bool   `json:"sell_by_cash" binding:"required"`
	SellAmountByCash int64  `json:"sell_amount_by_cash" binding:"required,gte=0"`
	SellExchangeRate int64  `json:"sell_exchange_rate" binding:"required,gte=0"`
	Comment          string `json:"comment" binding:"required"`
}

func (server *Server) createSellRequest(ctx *gin.Context) {
	var req createSellRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateSellRequestParams{
		SellTotalAmount:  req.SellTotalAmount,
		CurrencyFrom:     req.CurrencyFrom,
		CurrencyTo:       req.CurrencyTo,
		TgUsername:       req.TgUsername,
		SellByCard:       util.ToPgBool(req.SellByCard),
		SellAmountByCard: util.ToPgInt(req.SellAmountByCard),
		SellByCash:       util.ToPgBool(req.SellByCash),
		SellAmountByCash: util.ToPgInt(req.SellAmountByCash),
		SellExchangeRate: util.ToPgInt(req.SellExchangeRate),
		Comment:          req.Comment,
	}
	sellRequest, err := server.store.CreateSellRequest(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, sellRequest)
}

type getSellRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getSellRequest(ctx *gin.Context) {
	var req getSellRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	sellRequest, err := server.store.GetSellRequestById(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, sellRequest)
}

type listSellRequest struct {
	Limit  int32 `form:"limit" binding:"required,min=1"`
	Offset int32 `form:"offset" binding:"required,min=1"`
}

func (server *Server) listSellRequest(ctx *gin.Context) {
	var req listSellRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.ListSellRequestsParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}
	sellRequests, err := server.store.ListSellRequests(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, sellRequests)
}

type updateSellRequestUri struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type updateSellRequestJson struct {
	SellTotalAmount  *int64  `json:"sell_total_amount" binding:"omitempty,min=1"`
	CurrencyFrom     *string `json:"currency_from" binding:"omitempty,currency"`
	CurrencyTo       *string `json:"currency_to" binding:"omitempty,currency"`
	TgUsername       *string `json:"tg_username"`
	SellByCard       *bool   `json:"sell_by_card"`
	SellAmountByCard *int64  `json:"sell_amount_by_card" binding:"omitempty,gte=0"`
	SellByCash       *bool   `json:"sell_by_cash"`
	SellAmountByCash *int64  `json:"sell_amount_by_cash" binding:"omitempty,gte=0"`
	SellExchangeRate *int64  `json:"sell_exchange_rate" binding:"omitempty,min=1"`
	Comment          *string `json:"comment"`
}

func (server *Server) updateSellRequest(ctx *gin.Context) {
	var reqUri updateSellRequestUri
	var reqJson updateSellRequestJson
	err := ctx.ShouldBindUri(&reqUri)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = ctx.ShouldBindJSON(&reqJson)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.UpdateSellRequestParams{
		SellReqID: reqUri.ID,
		
		SellTotalAmount: pgtype.Int8{
			Int64: util.DerefInt64(reqJson.SellTotalAmount),
			Valid: reqJson.SellTotalAmount != nil,
		},
		CurrencyFrom: pgtype.Text{
			String: util.DerefStr(reqJson.CurrencyFrom),
			Valid:  reqJson.CurrencyFrom != nil,
		},
		CurrencyTo: pgtype.Text{
			String: util.DerefStr(reqJson.CurrencyTo),
			Valid:  reqJson.CurrencyTo != nil,
		},
		TgUsername: pgtype.Text{
			String: util.DerefStr(reqJson.TgUsername),
			Valid:  reqJson.TgUsername != nil,
		},
		SellByCard: pgtype.Bool{
			Bool:  util.DerefBool(reqJson.SellByCard),
			Valid: reqJson.SellByCard != nil,
		},
		SellAmountByCard: pgtype.Int8{
			Int64: util.DerefInt64(reqJson.SellAmountByCard),
			Valid: reqJson.SellAmountByCard != nil,
		},
		SellByCash: pgtype.Bool{
			Bool:  util.DerefBool(reqJson.SellByCash),
			Valid: reqJson.SellByCash != nil,
		},
		SellAmountByCash: pgtype.Int8{
			Int64: util.DerefInt64(reqJson.SellAmountByCash),
			Valid: reqJson.SellAmountByCash != nil,
		},
		SellExchangeRate: pgtype.Int8{
			Int64: util.DerefInt64(reqJson.SellExchangeRate),
			Valid: reqJson.SellExchangeRate != nil,
		},
		Comment: pgtype.Text{
			String: util.DerefStr(reqJson.Comment),
			Valid:  reqJson.Comment != nil,
		},
	}
	sellRequest, err := server.store.UpdateSellRequest(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, sellRequest)
}
