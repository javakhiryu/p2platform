package db

import (
	"context"
	"errors"
	"math"
	appErr "p2platform/errors"
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
				return appErr.ErrSellRequestNotFound
			}
			return appErr.ErrInternalServer
		}
		for _, sellRequest := range sellRequests {
			lockedAmounts, err := q.GetLockedAmountBySellRequest(ctx, sellRequest.SellReqID)
			if err != nil {
				if errors.Is(err, ErrNoRowsFound) {
					continue
				}
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
			results.SellRequests = append(results.SellRequests, GetSellRequestTxResult{
				SellRequest:        sellRequest,
				TotalLockedAmount:  totalLockedAmount,
				LockedAmountByCard: lockedAmountByCard,
				LockedAmountByCash: lockedAmountByCash,
			})
		}
		totalCount, err := q.CountOfSellRequestsByTelegramId(ctx, params.TelegramId)
		if err != nil {
			return appErr.ErrFailedToGetSellRequests
		}
		results.TotalPages = int32(math.Ceil(float64(totalCount) / float64(params.Limit)))
		return nil
	})
	return results, err
}
