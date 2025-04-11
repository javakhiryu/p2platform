package api

import (
	"net/http"
	db "p2platform/db/sqlc"
	"p2platform/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createBuyRequest struct {
	SellReqID       int32  `json:"sell_req_id" binding:"min=1"`
	BuyTotalAmount  int64  `json:"buy_total_amount" binding:"min=1"`
	TgUsername      string `json:"tg_username" binding:"required"`
	BuyAmountByCard int64  `json:"buy_amount_by_card" binding:"gte=0"`
	BuyAmountByCash int64  `json:"buy_amount_by_cash" binding:"min=0"`
}

func (server *Server) createBuyRequest(ctx *gin.Context) {
	var req createBuyRequest
	buyByCard := false
	buyByCash := false
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !(req.BuyAmountByCard+req.BuyAmountByCash == req.BuyTotalAmount) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "sum of amounts by cash and card is not equal to total buy amount",
		})
		return
	}
	if req.BuyAmountByCard > 0 {
		buyByCard = true
	}
	if req.BuyAmountByCash > 0 {
		buyByCash = true
	}

	arg := db.CreateBuyRequestTxParams{
		BuyReqID:        uuid.New(),
		SellReqID:       req.SellReqID,
		BuyTotalAmount:  req.BuyTotalAmount,
		TgUsername:      req.TgUsername,
		BuyByCard:       util.ToPgBool(buyByCard),
		BuyAmountByCard: util.ToPgInt(req.BuyAmountByCard),
		BuyByCash:       util.ToPgBool(buyByCash),
		BuyAmountByCash: util.ToPgInt(req.BuyAmountByCash),
	}
	buyRequest, err := server.store.CreateBuyRequestTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, buyRequest)
}

type getBuyRequestRequest struct {
	BuyRequestId string `uri:"id" binding:"required,uuid"`
}

func (server *Server) getBuyRequest(ctx *gin.Context) {
	var req getBuyRequestRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	uid, err := uuid.Parse(req.BuyRequestId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}
	buyRequest, err := server.store.GetBuyRequestById(ctx, uid)
	if err != nil {
		if err == db.ErrNoRowsFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, buyRequest)
}

type listBuyRequest struct {
	SellReqId int32 `form:"sell_req_id" binding:"min=1"`
	PageId    int32 `form:"page_id" binding:"min=1"`
	PageSize  int32 `form:"page_size" binding:"min=5,max=10"`
}

func (server *Server) listBuyRequests(ctx *gin.Context) {
	var req listBuyRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.ListBuyRequestsParams{
		SellReqID: req.SellReqId,
		Limit:     req.PageSize,
		Offset:    (req.PageId - 1) * req.PageSize,
	}
	buyRequests, err := server.store.ListBuyRequests(ctx, arg)
	if err != nil {
		if err == db.ErrNoRowsFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, buyRequests)
}

type closeBuyRequestRequest struct {
	BuyRequestId string `uri:"id" binding:"required,uuid"`
}

type closeBuyRequestResponse struct {
	CloseConfirmedBySeller bool       `json:"close_confirmed_by_seller"`
	SellerConfirmedAt      *time.Time `json:"seller_confirmed_at"`
	CloseConfirmedByBuyer  bool       `json:"close_confirmed_by_buyer"`
	BuyerConfirmedAt       *time.Time `json:"buyer_confirmed_at"`
	IsClosed               bool       `json:"is_closed"`
	ClosedAt               *time.Time `json:"closed_at"`
}

func (server *Server) closeBuyRequestBySeller(ctx *gin.Context) {
	var req closeBuyRequestRequest
	err := ctx.ShouldBindUri(&req)
	if err !=nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	uid, err := uuid.Parse(req.BuyRequestId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}
	arg := db.CloseBuyRequestTxParams{
		BuyRequestId: uid,
		IsSeller: true,
	}
	result, err := server.store.CloseBuyRequestTx(ctx, arg)
	if err !=nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	response :=closeBuyRequestResponse{
		CloseConfirmedBySeller: result.CloseConfirmedBySeller,
		SellerConfirmedAt: result.SellerConfirmedAt,
		CloseConfirmedByBuyer: result.CloseConfirmedByBuyer,
		BuyerConfirmedAt: result.BuyerConfirmedAt,
		IsClosed: result.IsClosed,
		ClosedAt: result.ClosedAt,
	}
	ctx.JSON(http.StatusOK, response)
}

func (server *Server) closeBuyRequestByBuyer(ctx *gin.Context) {
	var req closeBuyRequestRequest
	err := ctx.ShouldBindUri(&req)
	if err !=nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	uid, err := uuid.Parse(req.BuyRequestId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}
	arg := db.CloseBuyRequestTxParams{
		BuyRequestId: uid,
		IsSeller: false,
	}
	result, err := server.store.CloseBuyRequestTx(ctx, arg)
	if err !=nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	response :=closeBuyRequestResponse{
		CloseConfirmedBySeller: result.CloseConfirmedBySeller,
		SellerConfirmedAt: result.SellerConfirmedAt,
		CloseConfirmedByBuyer: result.CloseConfirmedByBuyer,
		BuyerConfirmedAt: result.BuyerConfirmedAt,
		IsClosed: result.IsClosed,
		ClosedAt: result.ClosedAt,
	}
	ctx.JSON(http.StatusOK, response)
}
