package db

import (
	"context"
	"errors"
	"fmt"
	appErr "p2platform/errors"
	"p2platform/util"
	"time"

	"github.com/google/uuid"
)

type CloseBuyRequestTxParams struct {
	BuyRequestId uuid.UUID
	IsSeller     bool
	IsBuyer      bool
}

type CloseBuyRequestTxResult struct {
	CloseConfirmedBySeller bool       `json:"close_confirmed_by_seller"`
	SellerConfirmedAt      *time.Time `json:"seller_confirmed_at"`
	CloseConfirmedByBuyer  bool       `json:"close_confirmed_by_buyer"`
	BuyerConfirmedAt       *time.Time `json:"buyer_confirmed_at"`
	BuyRequestState        string     `json:"buy_request_state"`
	StateUpdatedAt         *time.Time `json:"state_updated_at"`
}

func (store *SQLStore) CloseBuyRequestTx(ctx context.Context, arg CloseBuyRequestTxParams) (CloseBuyRequestTxResult, error) {
	var result CloseBuyRequestTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		buyRequest, err := store.GetBuyRequestById(ctx, arg.BuyRequestId)
		if err != nil {
			if errors.Is(err, ErrNoRowsFound) {
				return appErr.ErrBuyRequestsNotFound
			}
			return appErr.ErrFailedToGetBuyRequests
		}
		if buyRequest.State != "open" {
			return appErr.ErrBuyRequestAlreadyClosedOrExpired
		}
		if arg.IsSeller {
			sellerArgs := CloseConfirmBySellerParams{
				CloseConfirmBySeller: util.ToPgBool(true),
				BuyReqID:             arg.BuyRequestId,
			}
			err := store.CloseConfirmBySeller(ctx, sellerArgs)
			if err != nil {
				return appErr.ErrFailedToCloseBuyRequests
			}
		}
		if arg.IsBuyer {
			buyerArgs := CloseConfirmByBuyerParams{
				CloseConfirmByBuyer: util.ToPgBool(true),
				BuyReqID:            arg.BuyRequestId,
			}
			err := store.CloseConfirmByBuyer(ctx, buyerArgs)
			if err != nil {
				return appErr.ErrFailedToCloseBuyRequests
			}
		}
		buyRequest, err = store.GetBuyRequestById(ctx, arg.BuyRequestId)
		if err != nil {
			return appErr.ErrFailedToGetBuyRequests
		}

		var (
			sellerConfirmedAt = util.GetValidTime(buyRequest.SellerConfirmedAt)
			buyerConfirmedAt  = util.GetValidTime(buyRequest.BuyerConfirmedAt)
			stateChangedAt    = util.GetValidTime(buyRequest.StateUpdatedAt)
		)

		if buyRequest.CloseConfirmByBuyer.Bool && buyRequest.CloseConfirmBySeller.Bool {
			closeBuyRequestArgs := ChangeStateBuyRequestParams{
				State:    "closed",
				BuyReqID: arg.BuyRequestId,
			}
			closedBuyRequest, err := store.ChangeStateBuyRequest(ctx, closeBuyRequestArgs)
			if err != nil {
				return fmt.Errorf("Failed to close buy request")
			}
			result = CloseBuyRequestTxResult{
				CloseConfirmedBySeller: closedBuyRequest.CloseConfirmBySeller.Bool,
				SellerConfirmedAt:      &closedBuyRequest.SellerConfirmedAt.Time,
				CloseConfirmedByBuyer:  closedBuyRequest.CloseConfirmByBuyer.Bool,
				BuyerConfirmedAt:       &closedBuyRequest.BuyerConfirmedAt.Time,
				BuyRequestState:        closedBuyRequest.State,
				StateUpdatedAt:         &closedBuyRequest.StateUpdatedAt.Time,
			}
		} else {
			result = CloseBuyRequestTxResult{
				CloseConfirmedBySeller: buyRequest.CloseConfirmBySeller.Bool,
				SellerConfirmedAt:      sellerConfirmedAt,
				CloseConfirmedByBuyer:  buyRequest.CloseConfirmByBuyer.Bool,
				BuyerConfirmedAt:       buyerConfirmedAt,
				BuyRequestState:        buyRequest.State,
				StateUpdatedAt:         stateChangedAt,
			}
		}

		return nil
	})
	return result, err
}
