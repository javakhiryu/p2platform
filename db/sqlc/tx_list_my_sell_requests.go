package db

import (
	"context"
	"errors"
	"fmt"
)

type ListMySellRequestsTxParams struct {
	Limit      int32
	Offset     int32
	TelegramId int64
}

func (store *SQLStore) ListMySellRequeststTx(ctx context.Context, params ListMySellRequestsTxParams) (ListSellRequeststTxResults, error) {
	var results ListSellRequeststTxResults
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		arg := ListSellRequestsByTelegramIdParams{
			Limit:      params.Limit,
			Offset:     params.Offset,
			TelegramID: params.TelegramId,
		}
		sellRequests, err := q.ListSellRequestsByTelegramId(ctx, arg)
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
			results.SellRequests = append(results.SellRequests, GetSellRequestTxResult{
				SellRequest:        sellRequest,
				TotalLockedAmount:  totalLockedAmount,
				LockedAmountByCard: lockedAmountByCard,
				LockedAmountByCash: lockedAmountByCash,
			})
		}
		return nil
	})
	return results, err
}
