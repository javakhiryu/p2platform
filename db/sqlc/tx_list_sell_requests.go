package db

import (
	"context"
	"errors"
	"math"
	appErr "p2platform/errors"

	"github.com/google/uuid"
)

type ListSellRequeststTxParams struct {
	Limit   int32
	Offset  int32
	SpaceId uuid.UUID
}

type ListSellRequeststTxResults struct {
	SellRequests []GetSellRequestTxResult `json:"sell_requests"`
	TotalPages   int32                    `json:"total_pages"`
}

func (store *SQLStore) ListSellRequeststTx(ctx context.Context, params ListSellRequeststTxParams, requesterTelegramID int64) (ListSellRequeststTxResults, error) {
	var results ListSellRequeststTxResults
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		argIsUser := IsUserInSpaceParams{
			UserID:  requesterTelegramID,
			SpaceID: params.SpaceId,
		}

		isMember, err := q.IsUserInSpace(ctx, argIsUser)
		if err != nil {
			return appErr.ErrForbidden
		}
		if !isMember {
			return appErr.ErrForbidden
		}

		arg := ListSellRequestsBySpaceParams{
			SpaceID: params.SpaceId,
			Limit:   params.Limit,
			Offset:  params.Offset,
		}
		sellRequests, err := q.ListSellRequestsBySpace(ctx, arg)
		if err != nil {
			if errors.Is(err, ErrNoRowsFound) {
				return appErr.ErrSellRequestNotFound
			}
			return appErr.ErrFailedToGetSellRequests
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

		totalCount, err := q.CountOfSellRequestsBySpace(ctx, arg.SpaceID)
		if err != nil {
			return appErr.ErrFailedToGetSellRequests
		}
		results.TotalPages = int32(math.Ceil(float64(totalCount) / float64(params.Limit)))
		return nil
	})
	return results, err
}
