package db

import (
	"context"
	"errors"
	appErr "p2platform/errors"
	"p2platform/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateBuyRequestTxParams struct {
	BuyReqID        uuid.UUID   `json:"buy_req_id"`
	SellReqID       int32       `json:"sell_req_id"`
	BuyTotalAmount  int64       `json:"buy_total_amount"`
	TelegramId      int64       `json:"telegram_id"`
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
		sellRequest, err := q.GetSellRequestForUpdate(ctx, arg.SellReqID)
		if err != nil {
			if errors.Is(err, ErrNoRowsFound) {
				return appErr.ErrSellRequestNotFound
			}
			return appErr.ErrFailedToGetSellRequests
		}
		lockedAmounts, err := q.GetLockedAmountBySellRequest(ctx, arg.SellReqID)
		if err != nil && !errors.Is(err, ErrNoRowsFound) {
			return appErr.ErrFailedToGetLockedAmountBySellRequest
		}

		var totalLockedAmount int64
		var lockedAmountByCard int64
		var lockedAmountByCash int64

		if !sellRequest.IsActual.Bool {
			return appErr.ErrSellRequestIsNotActual
		}
		if sellRequest.IsDeleted.Bool {
			return appErr.ErrSellRequestDeleted
		}

		for _, lockedAmount := range lockedAmounts {
			totalLockedAmount += lockedAmount.LockedTotalAmount
		}
		if sellRequest.SellTotalAmount-totalLockedAmount < arg.BuyTotalAmount {
			return appErr.ErrInsuficientTotalFunds
		}
		if arg.BuyByCard.Bool == true {
			for _, lockedAmount := range lockedAmounts {
				lockedAmountByCard += lockedAmount.LockedByCard.Int64
			}
			if sellRequest.SellAmountByCard.Int64-lockedAmountByCard < arg.BuyAmountByCard.Int64 {
				return appErr.ErrInsuficientCardFunds
			}
		}
		if arg.BuyByCash.Bool == true {
			for _, lockedAmount := range lockedAmounts {
				lockedAmountByCash += lockedAmount.LockedByCash.Int64
			}
			if sellRequest.SellAmountByCash.Int64-lockedAmountByCash < arg.BuyAmountByCash.Int64 {
				return appErr.ErrInsuficientCashFunds
			}
		}

		user, err := q.GetUser(ctx, arg.TelegramId)
		if err != nil {
			if errors.Is(err, ErrNoRowsFound) {
				return appErr.ErrUserNotFound
			}
			return appErr.ErrInternalServer
		}

		// 5. Создаём buy_request
		buyRequest, err := q.CreateBuyRequest(ctx, CreateBuyRequestParams{
			BuyReqID:        arg.BuyReqID,
			SellReqID:       arg.SellReqID,
			BuyTotalAmount:  arg.BuyTotalAmount,
			TelegramID:      arg.TelegramId,
			TgUsername:      user.TgUsername,
			BuyByCard:       arg.BuyByCard,
			BuyAmountByCard: arg.BuyAmountByCard,
			BuyByCash:       arg.BuyByCash,
			BuyAmountByCash: arg.BuyAmountByCash,
		})
		if err != nil {
			return appErr.ErrFailedToCreateBuyRequest
		}
		lockedAmount, err := q.CreateLockedAmount(ctx, CreateLockedAmountParams{
			SellReqID:         arg.SellReqID,
			BuyReqID:          arg.BuyReqID,
			LockedTotalAmount: arg.BuyTotalAmount,
			LockedByCard:      arg.BuyAmountByCard,
			LockedByCash:      arg.BuyAmountByCash,
		})
		if err != nil {
			return appErr.ErrFailedToCreateLockedAmount
		}

		remainingAmount := sellRequest.SellTotalAmount - (totalLockedAmount + arg.BuyTotalAmount)

		if remainingAmount == 0 {
			_, err := q.OpenCloseSellRequest(ctx, OpenCloseSellRequestParams{
				IsActual:  util.ToPgBool(false),
				SellReqID: arg.SellReqID,
			})
			if err != nil {
				return appErr.ErrFailedToCloseSellRequest
			}
		}

		sellRequest, err = q.GetSellRequestById(ctx, arg.SellReqID)
		if err != nil {
			return appErr.ErrFailedToGetSellRequests
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
