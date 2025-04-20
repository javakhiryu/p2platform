package db

import (
	"context"
	"errors"
	appErr "p2platform/errors"
	"github.com/jackc/pgx/v5/pgtype"
)


func (store *SQLStore) DeleteSellRequestTx(ctx context.Context, sellRequestId int32) (bool, error) {
	var result pgtype.Bool
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		sellRequest, err := q.GetSellRequestForUpdate(ctx, sellRequestId)
		if err != nil {
			if errors.Is(err, ErrNoRowsFound) {
				return appErr.ErrSellRequestNotFound
			}
			return appErr.ErrFailedToGetSellRequests
		}

		if sellRequest.IsDeleted.Bool {
			return appErr.ErrSellRequestAlreadyDeleted
		}

		result, err = q.DeleteSellRequest(ctx, sellRequestId)
		if err != nil {
			return appErr.ErrFailedToDeleteSellRequest
		}

		err = q.CloseBuyRequestBySellRequest(ctx, sellRequestId)
		if err != nil {
			return appErr.ErrFailedToCloseBuyRequests
		}

		return nil
	})
	return result.Bool, err
}
