package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)


func (store *SQLStore) DeleteSellRequestTx(ctx context.Context, sellRequestId int32) (bool, error) {
	var result pgtype.Bool
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		sellRequest, err := q.GetSellRequestForUpdate(ctx, sellRequestId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return err
			}
			return err
		}

		if sellRequest.IsDeleted.Bool {
			return ErrSellRequestAlreadyDeleted
		}

		result, err = q.DeleteSellRequest(ctx, sellRequestId)
		if err != nil {
			return fmt.Errorf("failed to delete sell request: %w", err)
		}

		err = q.CloseBuyRequestBySellRequest(ctx, sellRequestId)
		if err != nil {
			return fmt.Errorf("failed to close buy request(s): %w", err)
		}

		return nil
	})
	return result.Bool, err
}
