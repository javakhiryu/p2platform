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
}

type CloseBuyRequestTxResult struct {
	CloseConfirmedBySeller bool
	SellerConfirmedAt      *time.Time
	CloseConfirmedByBuyer  bool
	BuyerConfirmedAt       *time.Time
	IsClosed               bool
	ClosedAt               *time.Time
}

func (store *SQLStore) CloseBuyRequestTx(ctx context.Context, arg CloseBuyRequestTxParams) (CloseBuyRequestTxResult, error) {
	var result CloseBuyRequestTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		_, err := store.GetBuyRequestById(ctx, arg.BuyRequestId)
		if err != nil {
			if errors.Is(err, ErrNoRowsFound) {
				return appErr.ErrBuyRequestsNotFound
			}
			return appErr.ErrFailedToGetBuyRequests
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
		} else {
			buyerArgs := CloseConfirmByBuyerParams{
				CloseConfirmByBuyer: util.ToPgBool(true),
				BuyReqID:            arg.BuyRequestId,
			}
			err := store.CloseConfirmByBuyer(ctx, buyerArgs)
			if err != nil {
				return appErr.ErrFailedToCloseBuyRequests
			}
		}
		buyRequest, err := store.GetBuyRequestById(ctx, arg.BuyRequestId)
		if err != nil {
			return appErr.ErrFailedToGetBuyRequests
		}

		var (
			sellerConfirmedAt = util.GetValidTime(buyRequest.SellerConfirmedAt)
			buyerConfirmedAt  = util.GetValidTime(buyRequest.BuyerConfirmedAt)
			closedAt          = util.GetValidTime(buyRequest.ClosedAt)
		)

		if buyRequest.CloseConfirmByBuyer.Bool && buyRequest.CloseConfirmBySeller.Bool {
			closeBuyRequestArgs := OpenCloseBuyRequestParams{
				IsClosed: util.ToPgBool(true),
				BuyReqID: arg.BuyRequestId,
			}
			closedBuyRequest, err := store.OpenCloseBuyRequest(ctx, closeBuyRequestArgs)
			if err != nil {
				return fmt.Errorf("Failed to close buy request")
			}
			result = CloseBuyRequestTxResult{
				CloseConfirmedBySeller: closedBuyRequest.CloseConfirmBySeller.Bool,
				SellerConfirmedAt:      &closedBuyRequest.SellerConfirmedAt.Time,
				CloseConfirmedByBuyer:  closedBuyRequest.CloseConfirmByBuyer.Bool,
				BuyerConfirmedAt:       &closedBuyRequest.BuyerConfirmedAt.Time,
				IsClosed:               closedBuyRequest.IsClosed.Bool,
				ClosedAt:               &closedBuyRequest.ClosedAt.Time,
			}
			err = store.ReleaseLockedAmountByBuyRequest(ctx, closeBuyRequestArgs.BuyReqID)
			if err != nil {
				return fmt.Errorf("Failed to release locked amount")
			}
		} else {
			result = CloseBuyRequestTxResult{
				CloseConfirmedBySeller: buyRequest.CloseConfirmBySeller.Bool,
				SellerConfirmedAt:      sellerConfirmedAt,
				CloseConfirmedByBuyer:  buyRequest.CloseConfirmByBuyer.Bool,
				BuyerConfirmedAt:       buyerConfirmedAt,
				IsClosed:               buyRequest.IsClosed.Bool,
				ClosedAt:               closedAt,
			}
		}

		return nil
	})
	return result, err
}
