package db

import (
	"context"
	"errors"
	"fmt"
)

type ListSellRequeststTxParams struct {
	Limit  int32
	Offset int32
}

type ListSellRequeststTxResults struct {
	SellRequests []GetSellRequestTxResult `json:"sell_requests"`
}

func (store *SQLStore) ListSellRequeststTx(ctx context.Context, params ListSellRequeststTxParams) (ListSellRequeststTxResults, error) {
	var results ListSellRequeststTxResults
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		arg := ListSellRequestsParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		}
		sellRequests, err := q.ListSellRequests(ctx, arg)
		if err != nil {
			if errors.Is(err, ErrNoRowsFound) {
				return fmt.Errorf(ErrSellRequestNotFound.Error(), err)
			}
			return fmt.Errorf("failed to list sell requests: %w", err)
		}
		for _, sellRequest := range sellRequests {
			lockedAmounts, err := q.GetLockedAmountBySellRequest(ctx, sellRequest.SellReqID)
			if err != nil {
				return fmt.Errorf("failed to get locked amounts for sell_request_id %d: %w", sellRequest.SellReqID, err)
			}

			var totalLockedAmount int64
			var lockedAmountByCard int64
			var lockedAmountByCash int64
			for _, lockedAmount := range lockedAmounts {
				totalLockedAmount += lockedAmount.LockedTotalAmount
				lockedAmountByCard += lockedAmount.LockedByCard.Int64
				lockedAmountByCash += lockedAmount.LockedByCash.Int64
			}
			result := GetSellRequestTxResult{
				SellRequest:        sellRequest,
				TotalLockedAmount:  totalLockedAmount,
				LockedAmountByCard: lockedAmountByCard,
				LockedAmountByCash: lockedAmountByCash,
			}

			results.SellRequests = append(results.SellRequests, result)
		}
		return nil
	})
	return results, err
}
