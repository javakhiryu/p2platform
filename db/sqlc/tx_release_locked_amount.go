package db

import (
	"context"
	"errors"
	"fmt"
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
		// Закрываем заявку всегда
		buyRequest, err := q.OpenCloseBuyRequest(ctx, OpenCloseBuyRequestParams{
			IsClosed: util.ToPgBool(true),
			BuyReqID: buyReqID,
		})
		if err != nil {
			return fmt.Errorf("failed to close buy request: %w", err)
		}

		result.BuyRequestID = buyRequest.BuyReqID
		result.IsBuyRequestClosed = buyRequest.IsClosed.Bool
		if buyRequest.ClosedAt.Valid {
			result.BuyRequestClosedAt = buyRequest.ClosedAt.Time
		}

		// Пробуем освободить заблокированную сумму, но это не обязательно
		err = q.ReleaseLockedAmountByBuyRequest(ctx, buyReqID)
		if err != nil {
			if errors.Is(err, ErrNoRowsFound) {
				// Не найдено — просто логируем, не прерываем транзакцию
				return nil
			}
			return fmt.Errorf("failed to release locked amount: %w", err)
		}

		lockedAmount, err := q.GetLockedAmount(ctx, buyReqID)
		if err != nil {
			if errors.Is(err, ErrNoRowsFound) {
				// Не найдено — допустимо
				return nil
			}
			return fmt.Errorf("failed to get locked amount: %w", err)
		}

		result.LockedAmountID = lockedAmount.ID
		result.IsLockedAmountReleased = lockedAmount.IsReleased.Bool
		if lockedAmount.ReleasedAt.Valid {
			result.LockedAmountReleasedAt = lockedAmount.ReleasedAt.Time
		}

		return nil
	})
	return result, err
}
