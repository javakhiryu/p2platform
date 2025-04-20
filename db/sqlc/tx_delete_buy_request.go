package db

import (
	"context"
	"errors"
	appErr "p2platform/errors"
	"p2platform/util"

	"github.com/google/uuid"
)

func (store *SQLStore) DeleteBuyRequestTx(ctx context.Context, buyReqID uuid.UUID) (bool, error) {
	var result bool
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		var buyRequest BuyRequest

		buyRequest, err = q.GetBuyRequestById(ctx, buyReqID)
		if err != nil {
			if errors.Is(err, ErrNoRowsFound) {
				return appErr.ErrBuyRequestsNotFound
			}
			return appErr.ErrFailedToGetBuyRequests
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

		err = q.DeleteBuyRequest(ctx, buyReqID)
		if err != nil {
			return appErr.ErrFailedToDeleteBuyRequest
		}
		result = true
		return nil
	})
	return result, err
}
