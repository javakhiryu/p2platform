package api

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	db "p2platform/db/sqlc"
	appErr "p2platform/errors"
	"p2platform/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type createBuyRequest struct {
	SellReqID       int32 `json:"sell_req_id" binding:"min=1"`
	BuyTotalAmount  int64 `json:"buy_total_amount" binding:"min=1"`
	BuyAmountByCard int64 `json:"buy_amount_by_card" binding:"gte=0"`
	BuyAmountByCash int64 `json:"buy_amount_by_cash" binding:"min=0"`
}

func (server *Server) createBuyRequest(ctx *gin.Context) {
	var req createBuyRequest
	buyByCard := false
	buyByCash := false
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidPayload.Status, ErrorResponse(appErr.ErrInvalidPayload))
		return
	}

	if err != nil {
		ctx.JSON(appErr.ErrInvalidUUID.Status, ErrorResponse(appErr.ErrInvalidUUID))
		return
	}

	telegramId, ok := GetTelegramIDFromContext(ctx)
	if !ok {
		return
	}

	_, err = server.store.GetSellRequestTx(ctx, req.SellReqID, telegramId)
	if err != nil {
		HandleAppError(ctx, err)
	}

	if !(req.BuyAmountByCard+req.BuyAmountByCash == req.BuyTotalAmount) {
		ctx.JSON(appErr.ErrTotalBuyAmountMismatch.Status, ErrorResponse(appErr.ErrTotalBuyAmountMismatch))
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
		TelegramId:      telegramId,
		BuyByCard:       util.ToPgBool(buyByCard),
		BuyAmountByCard: util.ToPgInt(req.BuyAmountByCard),
		BuyByCash:       util.ToPgBool(buyByCash),
		BuyAmountByCash: util.ToPgInt(req.BuyAmountByCash),
	}
	buyRequest, err := server.store.CreateBuyRequestTx(ctx, arg)
	if err != nil {
		HandleAppError(ctx, err)
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
		ctx.JSON(appErr.ErrInvalidUri.Status, ErrorResponse(appErr.ErrInvalidUri))
		return
	}
	uid, err := uuid.Parse(req.BuyRequestId)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidUUID.Status, ErrorResponse(appErr.ErrInvalidUUID))
		return
	}
	telegramId, exists := GetTelegramIDFromContext(ctx)
	if !exists {
		return
	}

	buyRequest, err := server.store.GetBuyRequestById(ctx, uid)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsFound) {
			ctx.JSON(appErr.ErrBuyRequestsNotFound.Status, ErrorResponse(appErr.ErrBuyRequestsNotFound))
			return
		}
		ctx.JSON(appErr.ErrFailedToGetBuyRequests.Status, ErrorResponse(appErr.ErrFailedToGetBuyRequests))
		return
	}
	arg := db.IsUserInSameSpaceAsSellerParams{
		UserID:   telegramId,
		UserID_2: buyRequest.TelegramID,
	}
	isSameSpace, err := server.store.IsUserInSameSpaceAsSeller(ctx, arg)
	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
		return
	}
	if !isSameSpace {
		ctx.JSON(appErr.ErrForbidden.Status, ErrorResponse(appErr.ErrForbidden))
		return
	}
	ctx.JSON(http.StatusOK, buyRequest)
}

type listBuyRequests struct {
	SellReqId int32 `form:"sell_req_id" binding:"min=1"`
	PageId    int32 `form:"page_id" binding:"min=1"`
	PageSize  int32 `form:"page_size" binding:"min=5,max=10"`
}

type listBuyRequestsResponse struct {
	BuyRequests []db.BuyRequest `json:"buy_requests"`
	TotalPages  int32           `json:"total_pages"`
}

func (server *Server) listBuyRequests(ctx *gin.Context) {
	var req listBuyRequests
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidQuery.Status, ErrorResponse(appErr.ErrInvalidQuery))
		return
	}
	telegramId, exists := GetTelegramIDFromContext(ctx)
	if !exists {
		return
	}
	_, err = server.store.GetSellRequestTx(ctx, req.SellReqId, telegramId)
	if err != nil {
		HandleAppError(ctx, err)
	}
	arg := db.ListBuyRequestsParams{
		SellReqID: req.SellReqId,
		Limit:     req.PageSize,
		Offset:    (req.PageId - 1) * req.PageSize,
	}
	buyRequests, err := server.store.ListBuyRequests(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsFound) {
			ctx.JSON(appErr.ErrBuyRequestsNotFound.Status, ErrorResponse(appErr.ErrBuyRequestsNotFound))
			return
		}
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
		return
	}
	totalBuyRequests, err := server.store.CountOfBuyRequests(ctx, req.SellReqId)
	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
	}

	totalPages := int32(math.Ceil(float64(totalBuyRequests) / float64(req.PageSize)))

	result := listBuyRequestsResponse{
		BuyRequests: buyRequests,
		TotalPages:  totalPages,
	}

	ctx.JSON(http.StatusOK, result)
}

type listMyBuyRequests struct {
	PageId   int32 `form:"page_id" binding:"min=1"`
	PageSize int32 `form:"page_size" binding:"min=5,max=10"`
	SpaceId  string `form:"space_id" binding:"required,uuid"`
}

func (server *Server) listMyBuyRequests(ctx *gin.Context) {
	var req listMyBuyRequests
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidQuery.Status, ErrorResponse(appErr.ErrInvalidQuery))
		return
	}
	telegramId, ok := GetTelegramIDFromContext(ctx)
	if !ok {
		return
	}

	uid, err := uuid.Parse(req.SpaceId)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidUUID.Status, ErrorResponse(appErr.ErrInvalidUUID))
		return
	}

	arg := db.ListBuyRequestsByUserInSpaceParams{
		SpaceID:    uid,
		UserID: telegramId,
		Limit:      req.PageSize,
		Offset:     (req.PageId - 1) * req.PageSize,
	}
	buyRequests, err := server.store.ListBuyRequestsByUserInSpace(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsFound) {
			ctx.JSON(appErr.ErrBuyRequestsNotFound.Status, ErrorResponse(appErr.ErrBuyRequestsNotFound))
			return
		}
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
		return
	}
	totalBuyRequests, err := server.store.CountBuyRequestsByUserInSpace(ctx, db.CountBuyRequestsByUserInSpaceParams{
		SpaceID:    uid,
		UserID: telegramId,
	})
	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
	}

	totalPages := int32(math.Ceil(float64(totalBuyRequests) / float64(req.PageSize)))

	result := listBuyRequestsResponse{
		BuyRequests: buyRequests,
		TotalPages:  totalPages,
	}

	ctx.JSON(http.StatusOK, result)
}

type closeBuyRequestRequest struct {
	BuyRequestId string `uri:"id" binding:"required,uuid"`
}

type closeBuyRequestResponse struct {
	CloseConfirmedBySeller bool       `json:"close_confirmed_by_seller"`
	SellerConfirmedAt      *time.Time `json:"seller_confirmed_at"`
	CloseConfirmedByBuyer  bool       `json:"close_confirmed_by_buyer"`
	BuyerConfirmedAt       *time.Time `json:"buyer_confirmed_at"`
	BuyRequestState        string     `json:"buy_request_state"`
	BuyRequestClosedAt     *time.Time `json:"state_updated_at"`
}

func (server *Server) closeBuyRequestBySeller(ctx *gin.Context) {
	var req closeBuyRequestRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidUri.Status, ErrorResponse(appErr.ErrInvalidUri))
		return
	}

	uid, err := uuid.Parse(req.BuyRequestId)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidUUID.Status, ErrorResponse(appErr.ErrInvalidUUID))
		return
	}

	telegramId, ok := GetTelegramIDFromContext(ctx)
	if !ok {
		return
	}

	buyRequest, err := server.store.GetBuyRequestById(ctx, uid)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsFound) {
			ctx.JSON(appErr.ErrBuyRequestsNotFound.Status, ErrorResponse(appErr.ErrBuyRequestsNotFound))
			return
		}
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
		return
	}
	sellRequest, err := server.store.GetSellRequestById(ctx, buyRequest.SellReqID)
	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
	}

	if telegramId != sellRequest.TelegramID {
		ctx.JSON(appErr.ErrNotSellRequestOwner.Status, ErrorResponse(appErr.ErrNotSellRequestOwner))
		return
	}

	arg := db.CloseBuyRequestTxParams{
		BuyRequestId: uid,
		IsSeller:     true,
	}
	result, err := server.store.CloseBuyRequestTx(ctx, arg)
	if err != nil {
		HandleAppError(ctx, err)
		return
	}
	response := closeBuyRequestResponse{
		CloseConfirmedBySeller: result.CloseConfirmedBySeller,
		SellerConfirmedAt:      result.SellerConfirmedAt,
		CloseConfirmedByBuyer:  result.CloseConfirmedByBuyer,
		BuyerConfirmedAt:       result.BuyerConfirmedAt,
		BuyRequestState:        result.BuyRequestState,
		BuyRequestClosedAt:     result.StateUpdatedAt,
	}
	ctx.JSON(http.StatusOK, response)
}

func (server *Server) closeBuyRequestByBuyer(ctx *gin.Context) {
	var req closeBuyRequestRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
	uid, err := uuid.Parse(req.BuyRequestId)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidUUID.Status, ErrorResponse(appErr.ErrInvalidUUID))
		return
	}

	telegramId, ok := GetTelegramIDFromContext(ctx)
	if !ok {
		return
	}

	buyRequest, err := server.store.GetBuyRequestById(ctx, uid)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsFound) {
			ctx.JSON(appErr.ErrBuyRequestsNotFound.Status, ErrorResponse(appErr.ErrBuyRequestsNotFound))
			return
		}
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
		return
	}
	if telegramId != buyRequest.TelegramID {
		ctx.JSON(appErr.ErrNotBuyRequestOwner.Status, ErrorResponse(appErr.ErrNotBuyRequestOwner))
		return
	}

	arg := db.CloseBuyRequestTxParams{
		BuyRequestId: uid,
		IsSeller:     false,
	}
	result, err := server.store.CloseBuyRequestTx(ctx, arg)
	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
	}
	response := closeBuyRequestResponse{
		CloseConfirmedBySeller: result.CloseConfirmedBySeller,
		SellerConfirmedAt:      result.SellerConfirmedAt,
		CloseConfirmedByBuyer:  result.CloseConfirmedByBuyer,
		BuyerConfirmedAt:       result.BuyerConfirmedAt,
		BuyRequestState:        result.BuyRequestState,
		BuyRequestClosedAt:     result.StateUpdatedAt,
	}
	ctx.JSON(http.StatusOK, response)
}

func (server *Server) CloseBuyRequestSellerBuyer(ctx *gin.Context) {
	var req closeBuyRequestRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
	uid, err := uuid.Parse(req.BuyRequestId)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidUUID.Status, ErrorResponse(appErr.ErrInvalidUUID))
		return
	}

	telegramId, ok := GetTelegramIDFromContext(ctx)
	if !ok {
		return
	}
	buyRequest, err := server.store.GetBuyRequestById(ctx, uid)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsFound) {
			ctx.JSON(appErr.ErrBuyRequestsNotFound.Status, ErrorResponse(appErr.ErrBuyRequestsNotFound))
			return
		}
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
		return
	}
	sellRequest, err := server.store.GetSellRequestById(ctx, buyRequest.SellReqID)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsFound) {
			ctx.JSON(appErr.ErrSellRequestNotFound.Status, ErrorResponse(appErr.ErrSellRequestNotFound))
			return
		}
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
		return
	}
	isSeller := false
	isBuyer := false
	if telegramId == buyRequest.TelegramID {
		isBuyer = true
	}
	if telegramId == sellRequest.TelegramID {
		isSeller = true
	}
	if !isSeller && !isBuyer {
		ctx.JSON(appErr.ErrNotBuyRequestOwner.Status, ErrorResponse(appErr.ErrNotBuyRequestOwner))
		ctx.JSON(appErr.ErrNotSellRequestOwner.Status, ErrorResponse(appErr.ErrNotSellRequestOwner))
		return
	}
	arg := db.CloseBuyRequestTxParams{
		BuyRequestId: uid,
		IsSeller:     isSeller,
		IsBuyer:      isBuyer,
	}
	result, err := server.store.CloseBuyRequestTx(ctx, arg)
	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
	}
	response := closeBuyRequestResponse{
		CloseConfirmedBySeller: result.CloseConfirmedBySeller,
		SellerConfirmedAt:      result.SellerConfirmedAt,
		CloseConfirmedByBuyer:  result.CloseConfirmedByBuyer,
		BuyerConfirmedAt:       result.BuyerConfirmedAt,
		BuyRequestState:        result.BuyRequestState,
		BuyRequestClosedAt:     result.StateUpdatedAt,
	}
	ctx.JSON(http.StatusOK, response)
}

type deleteBuyRequestRequest struct {
	BuyRequestId string `uri:"id" binding:"required,uuid"`
}

type deleteBuyRequestResponse struct {
	Message   string `json:"message"`
	IsDeleted bool   `json:"is_deleted"`
}

func (server *Server) DeleteBuyRequest(ctx *gin.Context) {
	var req closeBuyRequestRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(appErr.ErrInvalidUri.Status, ErrorResponse(appErr.ErrInvalidUri))
		return
	}

	uid, err := uuid.Parse(req.BuyRequestId)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidUUID.Status, ErrorResponse(appErr.ErrInvalidUUID))
		return
	}

	telegramId, ok := GetTelegramIDFromContext(ctx)
	if !ok {
		return
	}

	buyRequest, err := server.store.GetBuyRequestById(ctx, uid)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsFound) {
			ctx.JSON(appErr.ErrBuyRequestsNotFound.Status, ErrorResponse(appErr.ErrBuyRequestsNotFound))
			return
		}
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
		return
	}

	if telegramId != buyRequest.TelegramID {
		ctx.JSON(appErr.ErrNotBuyRequestOwner.Status, ErrorResponse(appErr.ErrNotBuyRequestOwner))
		return
	}

	isDeleted, err := server.store.DeleteBuyRequestTx(ctx, uid)
	if err != nil {
		HandleAppError(ctx, err)
		return
	}

	result := deleteBuyRequestResponse{
		Message:   fmt.Sprintf("Buy request with UID %s deleted successfully", uid),
		IsDeleted: isDeleted,
	}

	ctx.JSON(http.StatusOK, result)
}
