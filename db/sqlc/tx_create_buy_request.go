package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"p2platform/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateBuyRequestTxParams struct {
	BuyReqID        uuid.UUID   `json:"buy_req_id"`
	SellReqID       int32       `json:"sell_req_id"`
	BuyTotalAmount  int64       `json:"buy_total_amount"`
	TgUsername      string      `json:"tg_username"`
	BuyByCard       pgtype.Bool `json:"buy_by_card"`
	BuyAmountByCard pgtype.Int8 `json:"buy_amount_by_card"`
	BuyByCash       pgtype.Bool `json:"buy_by_cash"`
	BuyAmountByCash pgtype.Int8 `json:"buy_amount_by_cash"`
}

type CreateBuyRequestTxResult struct {
	BuyRequest   BuyRequest   `json:"buy_request"`
	LockedAmount LockedAmount `json:"locked_amount"`
	SellRequest  SellRequest  `json:"sell_request"`
}

func (store *SQLStore) CreateBuyRequestTx(ctx context.Context, arg CreateBuyRequestTxParams) (CreateBuyRequestTxResult, error) {
	var result CreateBuyRequestTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// 1. Получаем sell_request
		sellRequest, err := q.GetSellRequestForUpdate(ctx, arg.SellReqID)
		if err != nil {
			if errors.Is(err, ErrNoRowsFound) {
				return fmt.Errorf("sell request not found: %w", err)
			}
			return fmt.Errorf("failed to get sell request: %w", err)
		}
		// 2. Получаем все блокировки по sell_request
		lockedAmounts, err := q.GetLockedAmountBySellRequest(ctx, arg.SellReqID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to get locked amounts: %v", err)
		}

		var totalLockedAmount int64
		var lockedAmountByCard int64
		var lockedAmountByCash int64

		if !sellRequest.IsActual.Bool {
			return fmt.Errorf("sell request is not actual")
		}
		if sellRequest.IsDeleted.Bool {
			return fmt.Errorf("Sell request has been deleted")
		}

		for _, lockedAmount := range lockedAmounts {
			totalLockedAmount += lockedAmount.LockedTotalAmount
		}
		if sellRequest.SellTotalAmount-totalLockedAmount < arg.BuyTotalAmount {
			return fmt.Errorf("insufficient funds: not enough available amount in sell request")
		}
		if arg.BuyByCard.Bool == true {
			for _, lockedAmount := range lockedAmounts {
				lockedAmountByCard += lockedAmount.LockedByCard.Int64
			}
			if sellRequest.SellAmountByCard.Int64-lockedAmountByCard < arg.BuyAmountByCard.Int64 {
				return fmt.Errorf("insufficient funds: not enough available amount by card in sell request")
			}
		}
		if arg.BuyByCash.Bool == true {
			for _, lockedAmount := range lockedAmounts {
				lockedAmountByCash += lockedAmount.LockedByCash.Int64
			}
			if sellRequest.SellAmountByCash.Int64-lockedAmountByCash < arg.BuyAmountByCard.Int64 {
				return fmt.Errorf("insufficient funds: not enough available amount by cash in sell request")
			}
		}

		// 5. Создаём buy_request
		buyRequest, err := q.CreateBuyRequest(ctx, CreateBuyRequestParams{
			BuyReqID:        arg.BuyReqID,
			SellReqID:       arg.SellReqID,
			BuyTotalAmount:  arg.BuyTotalAmount,
			TgUsername:      arg.TgUsername,
			BuyByCard:       arg.BuyByCard,
			BuyAmountByCard: arg.BuyAmountByCard,
			BuyByCash:       arg.BuyByCash,
			BuyAmountByCash: arg.BuyAmountByCash,
		})
		if err != nil {
			return fmt.Errorf("failed to create buy request: %v", err)
		}
		lockedAmount, err := q.CreateLockedAmount(ctx, CreateLockedAmountParams{
			SellReqID:         arg.SellReqID,
			BuyReqID:          arg.BuyReqID,
			LockedTotalAmount: arg.BuyTotalAmount,
			LockedByCard:      arg.BuyAmountByCard,
			LockedByCash:      arg.BuyAmountByCash,
		})
		if err != nil {
			return fmt.Errorf("failed to create locked amount: %v", err)
		}

		remainingAmount := sellRequest.SellTotalAmount - (totalLockedAmount + arg.BuyTotalAmount)

		if remainingAmount == 0 {
			_, err := q.OpenCloseSellRequest(ctx, OpenCloseSellRequestParams{
				IsActual:  util.ToPgBool(false),
				SellReqID: arg.SellReqID,
			})
			if err != nil {
				return fmt.Errorf("failed to close sell request: %v", err)
			}
		}

		sellRequest, err = q.GetSellRequestById(ctx, arg.SellReqID)
		if err != nil {
			return fmt.Errorf("failed to get sell request: %v", err)
		}

		result = CreateBuyRequestTxResult{
			BuyRequest:   buyRequest,
			LockedAmount: lockedAmount,
			SellRequest:  sellRequest,
		}
		return nil
	})
	return result, err
}
