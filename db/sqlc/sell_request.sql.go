// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: sell_request.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createSellRequest = `-- name: CreateSellRequest :one
INSERT INTO sell_requests (
  sell_total_amount,
  sell_money_source,
  currency_from,
  currency_to,
  telegram_id,
  tg_username,
  sell_by_card,
  sell_amount_by_card,
  sell_by_cash,
  sell_amount_by_cash,
  sell_exchange_rate,
  comment
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING sell_req_id, sell_total_amount, sell_money_source, currency_from, currency_to, telegram_id, tg_username, sell_by_card, sell_amount_by_card, sell_by_cash, sell_amount_by_cash, sell_exchange_rate, is_actual, created_at, updated_at, is_deleted, comment
`

type CreateSellRequestParams struct {
	SellTotalAmount  int64       `json:"sell_total_amount"`
	SellMoneySource  string      `json:"sell_money_source"`
	CurrencyFrom     string      `json:"currency_from"`
	CurrencyTo       string      `json:"currency_to"`
	TelegramID       int64       `json:"telegram_id"`
	TgUsername       string      `json:"tg_username"`
	SellByCard       pgtype.Bool `json:"sell_by_card"`
	SellAmountByCard pgtype.Int8 `json:"sell_amount_by_card"`
	SellByCash       pgtype.Bool `json:"sell_by_cash"`
	SellAmountByCash pgtype.Int8 `json:"sell_amount_by_cash"`
	SellExchangeRate pgtype.Int8 `json:"sell_exchange_rate"`
	Comment          string      `json:"comment"`
}

func (q *Queries) CreateSellRequest(ctx context.Context, arg CreateSellRequestParams) (SellRequest, error) {
	row := q.db.QueryRow(ctx, createSellRequest,
		arg.SellTotalAmount,
		arg.SellMoneySource,
		arg.CurrencyFrom,
		arg.CurrencyTo,
		arg.TelegramID,
		arg.TgUsername,
		arg.SellByCard,
		arg.SellAmountByCard,
		arg.SellByCash,
		arg.SellAmountByCash,
		arg.SellExchangeRate,
		arg.Comment,
	)
	var i SellRequest
	err := row.Scan(
		&i.SellReqID,
		&i.SellTotalAmount,
		&i.SellMoneySource,
		&i.CurrencyFrom,
		&i.CurrencyTo,
		&i.TelegramID,
		&i.TgUsername,
		&i.SellByCard,
		&i.SellAmountByCard,
		&i.SellByCash,
		&i.SellAmountByCash,
		&i.SellExchangeRate,
		&i.IsActual,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.IsDeleted,
		&i.Comment,
	)
	return i, err
}

const deleteSellRequest = `-- name: DeleteSellRequest :one
UPDATE sell_requests
SET
  is_deleted = true,
  updated_at = now()
WHERE
  sell_req_id = $1
RETURNING is_deleted
`

func (q *Queries) DeleteSellRequest(ctx context.Context, sellReqID int32) (pgtype.Bool, error) {
	row := q.db.QueryRow(ctx, deleteSellRequest, sellReqID)
	var is_deleted pgtype.Bool
	err := row.Scan(&is_deleted)
	return is_deleted, err
}

const getSellRequestById = `-- name: GetSellRequestById :one
SELECT sell_req_id, sell_total_amount, sell_money_source, currency_from, currency_to, telegram_id, tg_username, sell_by_card, sell_amount_by_card, sell_by_cash, sell_amount_by_cash, sell_exchange_rate, is_actual, created_at, updated_at, is_deleted, comment FROM sell_requests WHERE sell_req_id = $1
`

func (q *Queries) GetSellRequestById(ctx context.Context, sellReqID int32) (SellRequest, error) {
	row := q.db.QueryRow(ctx, getSellRequestById, sellReqID)
	var i SellRequest
	err := row.Scan(
		&i.SellReqID,
		&i.SellTotalAmount,
		&i.SellMoneySource,
		&i.CurrencyFrom,
		&i.CurrencyTo,
		&i.TelegramID,
		&i.TgUsername,
		&i.SellByCard,
		&i.SellAmountByCard,
		&i.SellByCash,
		&i.SellAmountByCash,
		&i.SellExchangeRate,
		&i.IsActual,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.IsDeleted,
		&i.Comment,
	)
	return i, err
}

const getSellRequestForUpdate = `-- name: GetSellRequestForUpdate :one
SELECT sell_req_id, sell_total_amount, sell_money_source, currency_from, currency_to, telegram_id, tg_username, sell_by_card, sell_amount_by_card, sell_by_cash, sell_amount_by_cash, sell_exchange_rate, is_actual, created_at, updated_at, is_deleted, comment FROM sell_requests
WHERE sell_req_id = $1
FOR UPDATE
`

func (q *Queries) GetSellRequestForUpdate(ctx context.Context, sellReqID int32) (SellRequest, error) {
	row := q.db.QueryRow(ctx, getSellRequestForUpdate, sellReqID)
	var i SellRequest
	err := row.Scan(
		&i.SellReqID,
		&i.SellTotalAmount,
		&i.SellMoneySource,
		&i.CurrencyFrom,
		&i.CurrencyTo,
		&i.TelegramID,
		&i.TgUsername,
		&i.SellByCard,
		&i.SellAmountByCard,
		&i.SellByCash,
		&i.SellAmountByCash,
		&i.SellExchangeRate,
		&i.IsActual,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.IsDeleted,
		&i.Comment,
	)
	return i, err
}

const listSellRequests = `-- name: ListSellRequests :many
SELECT sell_req_id, sell_total_amount, sell_money_source, currency_from, currency_to, telegram_id, tg_username, sell_by_card, sell_amount_by_card, sell_by_cash, sell_amount_by_cash, sell_exchange_rate, is_actual, created_at, updated_at, is_deleted, comment FROM sell_requests
WHERE is_deleted = false
ORDER BY created_at ASC
LIMIT $1 
OFFSET $2
`

type ListSellRequestsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListSellRequests(ctx context.Context, arg ListSellRequestsParams) ([]SellRequest, error) {
	rows, err := q.db.Query(ctx, listSellRequests, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []SellRequest{}
	for rows.Next() {
		var i SellRequest
		if err := rows.Scan(
			&i.SellReqID,
			&i.SellTotalAmount,
			&i.SellMoneySource,
			&i.CurrencyFrom,
			&i.CurrencyTo,
			&i.TelegramID,
			&i.TgUsername,
			&i.SellByCard,
			&i.SellAmountByCard,
			&i.SellByCash,
			&i.SellAmountByCash,
			&i.SellExchangeRate,
			&i.IsActual,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.IsDeleted,
			&i.Comment,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSellRequestsByTelegramId = `-- name: ListSellRequestsByTelegramId :many
SELECT sell_req_id, sell_total_amount, sell_money_source, currency_from, currency_to, telegram_id, tg_username, sell_by_card, sell_amount_by_card, sell_by_cash, sell_amount_by_cash, sell_exchange_rate, is_actual, created_at, updated_at, is_deleted, comment FROM sell_requests
WHERE telegram_id = $1 AND is_deleted = false
ORDER BY created_at ASC
LIMIT $2 
OFFSET $3
`

type ListSellRequestsByTelegramIdParams struct {
	TelegramID int64 `json:"telegram_id"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
}

func (q *Queries) ListSellRequestsByTelegramId(ctx context.Context, arg ListSellRequestsByTelegramIdParams) ([]SellRequest, error) {
	rows, err := q.db.Query(ctx, listSellRequestsByTelegramId, arg.TelegramID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []SellRequest{}
	for rows.Next() {
		var i SellRequest
		if err := rows.Scan(
			&i.SellReqID,
			&i.SellTotalAmount,
			&i.SellMoneySource,
			&i.CurrencyFrom,
			&i.CurrencyTo,
			&i.TelegramID,
			&i.TgUsername,
			&i.SellByCard,
			&i.SellAmountByCard,
			&i.SellByCash,
			&i.SellAmountByCash,
			&i.SellExchangeRate,
			&i.IsActual,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.IsDeleted,
			&i.Comment,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const openCloseSellRequest = `-- name: OpenCloseSellRequest :one
UPDATE sell_requests
SET
  is_actual = $1,
  updated_at = now()
WHERE
  sell_req_id = $2
RETURNING sell_req_id, sell_total_amount, sell_money_source, currency_from, currency_to, telegram_id, tg_username, sell_by_card, sell_amount_by_card, sell_by_cash, sell_amount_by_cash, sell_exchange_rate, is_actual, created_at, updated_at, is_deleted, comment
`

type OpenCloseSellRequestParams struct {
	IsActual  pgtype.Bool `json:"is_actual"`
	SellReqID int32       `json:"sell_req_id"`
}

func (q *Queries) OpenCloseSellRequest(ctx context.Context, arg OpenCloseSellRequestParams) (SellRequest, error) {
	row := q.db.QueryRow(ctx, openCloseSellRequest, arg.IsActual, arg.SellReqID)
	var i SellRequest
	err := row.Scan(
		&i.SellReqID,
		&i.SellTotalAmount,
		&i.SellMoneySource,
		&i.CurrencyFrom,
		&i.CurrencyTo,
		&i.TelegramID,
		&i.TgUsername,
		&i.SellByCard,
		&i.SellAmountByCard,
		&i.SellByCash,
		&i.SellAmountByCash,
		&i.SellExchangeRate,
		&i.IsActual,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.IsDeleted,
		&i.Comment,
	)
	return i, err
}

const updateSellRequest = `-- name: UpdateSellRequest :one
UPDATE sell_requests
SET
    sell_total_amount = COALESCE($1, sell_total_amount),
    sell_money_source = COALESCE($2, sell_money_source),
    currency_from = COALESCE($3, currency_from),
    currency_to = COALESCE($4, currency_to),
    sell_by_card = COALESCE($5, sell_by_card),
    sell_amount_by_card = COALESCE($6, sell_amount_by_card),
    sell_by_cash = COALESCE($7, sell_by_cash),
    sell_amount_by_cash = COALESCE($8, sell_amount_by_cash),
    sell_exchange_rate = COALESCE($9, sell_exchange_rate),
    comment = COALESCE($10, comment),
    updated_at = CASE
        WHEN $1 IS NOT NULL
          OR $3 IS NOT NULL
          OR $4 IS NOT NULL
          OR $5 IS NOT NULL
          OR $6 IS NOT NULL
          OR $7 IS NOT NULL
          OR $8 IS NOT NULL
          OR $9 IS NOT NULL
          OR $10 IS NOT NULL
        THEN now()
        ELSE updated_at
    END
WHERE sell_req_id = $11
RETURNING sell_req_id, sell_total_amount, sell_money_source, currency_from, currency_to, telegram_id, tg_username, sell_by_card, sell_amount_by_card, sell_by_cash, sell_amount_by_cash, sell_exchange_rate, is_actual, created_at, updated_at, is_deleted, comment
`

type UpdateSellRequestParams struct {
	SellTotalAmount  pgtype.Int8 `json:"sell_total_amount"`
	SellMoneySource  pgtype.Text `json:"sell_money_source"`
	CurrencyFrom     pgtype.Text `json:"currency_from"`
	CurrencyTo       pgtype.Text `json:"currency_to"`
	SellByCard       pgtype.Bool `json:"sell_by_card"`
	SellAmountByCard pgtype.Int8 `json:"sell_amount_by_card"`
	SellByCash       pgtype.Bool `json:"sell_by_cash"`
	SellAmountByCash pgtype.Int8 `json:"sell_amount_by_cash"`
	SellExchangeRate pgtype.Int8 `json:"sell_exchange_rate"`
	Comment          pgtype.Text `json:"comment"`
	SellReqID        int32       `json:"sell_req_id"`
}

func (q *Queries) UpdateSellRequest(ctx context.Context, arg UpdateSellRequestParams) (SellRequest, error) {
	row := q.db.QueryRow(ctx, updateSellRequest,
		arg.SellTotalAmount,
		arg.SellMoneySource,
		arg.CurrencyFrom,
		arg.CurrencyTo,
		arg.SellByCard,
		arg.SellAmountByCard,
		arg.SellByCash,
		arg.SellAmountByCash,
		arg.SellExchangeRate,
		arg.Comment,
		arg.SellReqID,
	)
	var i SellRequest
	err := row.Scan(
		&i.SellReqID,
		&i.SellTotalAmount,
		&i.SellMoneySource,
		&i.CurrencyFrom,
		&i.CurrencyTo,
		&i.TelegramID,
		&i.TgUsername,
		&i.SellByCard,
		&i.SellAmountByCard,
		&i.SellByCash,
		&i.SellAmountByCash,
		&i.SellExchangeRate,
		&i.IsActual,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.IsDeleted,
		&i.Comment,
	)
	return i, err
}
