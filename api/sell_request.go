package api

import (
	"database/sql"
	"errors"
	"net/http"
	db "p2platform/db/sqlc"
	util "p2platform/util"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type createSellRequest struct {
	SellTotalAmount  int64  `json:"sell_total_amount" binding:"required,min=5"`
	SellMoneySource  string `json:"sell_money_source" binding:"required,source"`
	CurrencyFrom     string `json:"currency_from" binding:"required,currency"`
	CurrencyTo       string `json:"currency_to" binding:"required,currency"`
	TgUsername       string `json:"tg_username" binding:"required,max=100"`
	SellAmountByCard int64  `json:"sell_amount_by_card" binding:"gte=0"`
	SellAmountByCash int64  `json:"sell_amount_by_cash" binding:"gte=0"`
	SellExchangeRate int64  `json:"sell_exchange_rate" binding:"required,gt=0"`
	Comment          string `json:"comment" binding:"omitempty"`
}

func (server *Server) createSellRequest(ctx *gin.Context) {
	var req createSellRequest
	sellByCard := false
	sellByCash := false

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !(req.SellAmountByCard+req.SellAmountByCash == req.SellTotalAmount) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "sum of amounts by cash and card is not equal to total amount",
		})
		return
	}
	if req.SellAmountByCard > 0 {
		sellByCard = true
	}
	if req.SellAmountByCash > 0 {
		sellByCash = true
	}

	arg := db.CreateSellRequestParams{
		SellTotalAmount:  req.SellTotalAmount,
		SellMoneySource:  req.SellMoneySource,
		CurrencyFrom:     req.CurrencyFrom,
		CurrencyTo:       req.CurrencyTo,
		TgUsername:       req.TgUsername,
		SellByCard:       util.ToPgBool(sellByCard),
		SellAmountByCard: util.ToPgInt(req.SellAmountByCard),
		SellByCash:       util.ToPgBool(sellByCash),
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
		if err == db.ErrNoRowsFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if sellRequest.IsDeleted.Bool {
		ctx.JSON(http.StatusGone, gin.H{
			"error":      "sell request has been deleted",
			"deleted_at": sellRequest.UpdatedAt.UTC(),
		})
		return
	}
	ctx.JSON(http.StatusOK, sellRequest)
}

type listSellRequest struct {
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
	PageId   int32 `form:"page_id" binding:"required,min=1"`
}

func (server *Server) listSellRequest(ctx *gin.Context) {
	var req listSellRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.ListSellRequestsParams{
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}
	sellRequests, err := server.store.ListSellRequests(ctx, arg)
	if err != nil {
		if err == db.ErrNoRowsFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
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
	SellMoneySource  *string `json:"sell_money_source" binding:"omitempty,source"`
	CurrencyFrom     *string `json:"currency_from" binding:"omitempty,currency"`
	CurrencyTo       *string `json:"currency_to" binding:"omitempty,currency"`
	TgUsername       *string `json:"tg_username"`
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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = ctx.ShouldBindJSON(&reqJson)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	buyRequests, err := server.store.ListBuyRequests(ctx, db.ListBuyRequestsParams{
		SellReqID: reqUri.ID,
		Limit:     1,
		Offset:    0,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if len(buyRequests) > 0 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "sell request can't be updated because it has buy requests",
		})
		return
	}

	sellRequest, err := server.store.GetSellRequestById(ctx, reqUri.ID)
	if err != nil {
		if err == db.ErrNoRowsFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if sellRequest.IsDeleted.Bool {
		ctx.JSON(http.StatusGone, gin.H{
			"error":      "sell request has been deleted",
			"deleted_at": sellRequest.UpdatedAt.UTC(),
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "sum of amounts by cash and card is not equal to total amount",
		})
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
		TgUsername: pgtype.Text{
			String: util.DerefStr(reqJson.TgUsername),
			Valid:  reqJson.TgUsername != nil,
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	isDeleted, err := server.store.DeleteSellRequestTx(ctx, req.ID)
	if err != nil {
		switch {
		case errors.Is(err, db.ErrSellRequestAlreadyDeleted):
			ctx.JSON(http.StatusConflict, gin.H{"error": db.ErrSellRequestAlreadyDeleted.Error()})
			return
		case errors.Is(err, sql.ErrNoRows):
			ctx.JSON(http.StatusNotFound, gin.H{"error": db.ErrSEllRequestNotFound.Error()})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	result := deleteSellRequestResponse{
		IsDeleted: isDeleted,
	}
	ctx.JSON(http.StatusOK, result)
}
