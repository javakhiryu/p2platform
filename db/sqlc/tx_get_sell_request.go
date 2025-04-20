package db

import (
	"context"
	"database/sql"
	"errors"
	appErr "p2platform/errors"
)

type GetSellRequestTxResult struct {
	SellRequest        SellRequest `json:"sell_request"`
	TotalLockedAmount  int64       `json:"total_locked_amount"`
	LockedAmountByCard int64       `json:"locked_amount_by_card"`
	LockedAmountByCash int64       `json:"locked_amount_by_cash"`
}

func (store *SQLStore) GetSellRequestTx(ctx context.Context, sellReqID int32) (GetSellRequestTxResult, error) {
	var result GetSellRequestTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		sellRequest, err := q.GetSellRequestById(ctx, sellReqID)
		if err != nil {
			if errors.Is(err, ErrNoRowsFound) {
				return appErr.ErrSellRequestNotFound
			}
			return appErr.ErrFailedToGetSellRequests
		}
		result.SellRequest = sellRequest
		lockedAmounts, err := q.GetLockedAmountBySellRequest(ctx, sellReqID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return appErr.ErrFailedToGetLockedAmountBySellRequest
		}
		var totalLockedAmount int64
		var lockedAmountByCard int64
		var lockedAmountByCash int64
		for _, lockedAmount := range lockedAmounts {
			totalLockedAmount += lockedAmount.LockedTotalAmount
			lockedAmountByCard += lockedAmount.LockedByCard.Int64
			lockedAmountByCash += lockedAmount.LockedByCash.Int64
		}
		result.TotalLockedAmount = totalLockedAmount
		result.LockedAmountByCard = lockedAmountByCard
		result.LockedAmountByCash = lockedAmountByCash
		return nil
	})
	return result, err
}
