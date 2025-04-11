package db

import (
	"context"
	"errors"
	"fmt"
	"p2platform/util"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
				return fmt.Errorf("buy request not found: %w", err)
			}
			return fmt.Errorf("failed to get buy request: %w", err)
		}
		if arg.IsSeller {
			sellerArgs := CloseConfirmBySellerParams{
				CloseConfirmBySeller: util.ToPgBool(true),
				BuyReqID:             arg.BuyRequestId,
			}
			err := store.CloseConfirmBySeller(ctx, sellerArgs)
			if err != nil {
				return fmt.Errorf("Close confirmation by seller failed")
			}
		} else {
			buyerArgs := CloseConfirmByBuyerParams{
				CloseConfirmByBuyer: util.ToPgBool(true),
				BuyReqID:            arg.BuyRequestId,
			}
			err := store.CloseConfirmByBuyer(ctx, buyerArgs)
			if err != nil {
				return fmt.Errorf("Close confirmatin by buyer failed")
			}
		}
		buyRequest, err := store.GetBuyRequestById(ctx, arg.BuyRequestId)
		if err != nil {
			return fmt.Errorf("failed to re-fetch buy request: %w", err)
		}


		var (
			sellerConfirmedAt = getValidTime(buyRequest.SellerConfirmedAt)
			buyerConfirmedAt  = getValidTime(buyRequest.BuyerConfirmedAt)
			closedAt          = getValidTime(buyRequest.ClosedAt)
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

func getValidTime(t pgtype.Timestamptz) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}
