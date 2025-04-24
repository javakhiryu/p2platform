package db

import (
	"context"
	"errors"
	appErr "p2platform/errors"
	"p2platform/util"
	"time"

	"github.com/google/uuid"
)

type ReleaseLockedAmountTxResult struct {
	BuyRequestID           uuid.UUID `json:"buy_request_id"`
	IsBuyRequestClosed     bool      `json:"is_buy_request_closed"`
	BuyRequestClosedAt     time.Time `json:"buy_request_closed_at"`
	LockedAmountID         int32     `json:"locked_amount_id"`
	IsLockedAmountReleased bool      `json:"is_locked_amount_released"`
	LockedAmountReleasedAt time.Time `json:"locked_amount_released_at"`
}

func (store *SQLStore) ReleaseLockedAmountTx(ctx context.Context, buyReqID uuid.UUID) (result ReleaseLockedAmountTxResult, err error) {
	err = store.execTx(ctx, func(q *Queries) error {
		buyRequest, err := q.OpenCloseBuyRequest(ctx, OpenCloseBuyRequestParams{
			IsClosed: util.ToPgBool(true),
			BuyReqID: buyReqID,
		})
		if err != nil {
			return appErr.ErrFailedToCloseBuyRequests
		}

		result.BuyRequestID = buyRequest.BuyReqID
		result.IsBuyRequestClosed = buyRequest.IsClosed.Bool
		if buyRequest.ClosedAt.Valid {
			result.BuyRequestClosedAt = buyRequest.ClosedAt.Time
		}

		err = q.ReleaseLockedAmountByBuyRequest(ctx, buyReqID)
		if err != nil && !errors.Is(err, ErrNoRowsFound) {
			return appErr.ErrFailedToReleaseLockedAmount
		}

		lockedAmount, err := q.GetLockedAmount(ctx, buyReqID)
		if err != nil && !errors.Is(err, ErrNoRowsFound) {
			return appErr.ErrFailedToGetLockedAmountByBuyRequest
		}

		result.LockedAmountID = lockedAmount.ID
		result.IsLockedAmountReleased = lockedAmount.IsReleased.Bool
		if lockedAmount.ReleasedAt.Valid {
			result.LockedAmountReleasedAt = lockedAmount.ReleasedAt.Time
		}

		sellRequest, err := q.GetSellRequestForUpdate(ctx, buyRequest.SellReqID)
		if err != nil {
			return appErr.ErrFailedToGetSellRequests
		}

		if !sellRequest.IsDeleted.Bool && !sellRequest.IsActual.Bool {
			arg := OpenCloseSellRequestParams{
				IsActual:  util.ToPgBool(true),
				SellReqID: sellRequest.SellReqID,
			}
			_, err = q.OpenCloseSellRequest(ctx, arg)
			if err != nil {
				return appErr.ErrFailedToOpenSellRequest
			}
		}

		return nil
	})
	return result, err
}
