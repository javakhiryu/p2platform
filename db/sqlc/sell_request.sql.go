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
  sell_amount,
  currency_from,
  currency_to,
  tg_username,
  sell_by_card,
  sell_amount_by_card,
  sell_by_cash,
  sell_amount_by_cash,
  sell_exchange_rate,
  comment
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING sell_req_id, sell_amount, currency_from, currency_to, tg_username, sell_by_card, sell_amount_by_card, sell_by_cash, sell_amount_by_cash, sell_exchange_rate, is_actual, created_at, updated_at, is_deleted, comment
`

type CreateSellRequestParams struct {
	SellAmount       int64       `json:"sell_amount"`
	CurrencyFrom     string      `json:"currency_from"`
	CurrencyTo       string      `json:"currency_to"`
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
		arg.SellAmount,
		arg.CurrencyFrom,
		arg.CurrencyTo,
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
		&i.SellAmount,
		&i.CurrencyFrom,
		&i.CurrencyTo,
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
RETURNING sell_req_id, sell_amount, currency_from, currency_to, tg_username, sell_by_card, sell_amount_by_card, sell_by_cash, sell_amount_by_cash, sell_exchange_rate, is_actual, created_at, updated_at, is_deleted, comment
`

func (q *Queries) DeleteSellRequest(ctx context.Context, sellReqID int32) (SellRequest, error) {
	row := q.db.QueryRow(ctx, deleteSellRequest, sellReqID)
	var i SellRequest
	err := row.Scan(
		&i.SellReqID,
		&i.SellAmount,
		&i.CurrencyFrom,
		&i.CurrencyTo,
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

const getSellRequestById = `-- name: GetSellRequestById :one
SELECT sell_req_id, sell_amount, currency_from, currency_to, tg_username, sell_by_card, sell_amount_by_card, sell_by_cash, sell_amount_by_cash, sell_exchange_rate, is_actual, created_at, updated_at, is_deleted, comment FROM sell_requests WHERE sell_req_id = $1
`

func (q *Queries) GetSellRequestById(ctx context.Context, sellReqID int32) (SellRequest, error) {
	row := q.db.QueryRow(ctx, getSellRequestById, sellReqID)
	var i SellRequest
	err := row.Scan(
		&i.SellReqID,
		&i.SellAmount,
		&i.CurrencyFrom,
		&i.CurrencyTo,
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
SELECT sell_req_id, sell_amount, currency_from, currency_to, tg_username, sell_by_card, sell_amount_by_card, sell_by_cash, sell_amount_by_cash, sell_exchange_rate, is_actual, created_at, updated_at, is_deleted, comment FROM sell_requests
WHERE is_deleted = false
ORDER BY created_at DESC
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
			&i.SellAmount,
			&i.CurrencyFrom,
			&i.CurrencyTo,
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
RETURNING sell_req_id, sell_amount, currency_from, currency_to, tg_username, sell_by_card, sell_amount_by_card, sell_by_cash, sell_amount_by_cash, sell_exchange_rate, is_actual, created_at, updated_at, is_deleted, comment
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
		&i.SellAmount,
		&i.CurrencyFrom,
		&i.CurrencyTo,
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
    sell_amount = COALESCE($1, sell_amount),
    currency_from = COALESCE($2, currency_from),
    currency_to = COALESCE($3, currency_to),
    tg_username = COALESCE($4, tg_username),
    sell_by_card = COALESCE($5, sell_by_card),
    sell_amount_by_card = COALESCE($6, sell_amount_by_card),
    sell_by_cash = COALESCE($7, sell_by_cash),
    sell_amount_by_cash = COALESCE($8, sell_amount_by_cash),
    sell_exchange_rate = COALESCE($9, sell_exchange_rate),
    comment = COALESCE($10, comment),
    updated_at = CASE
        WHEN $1 IS NOT NULL
          OR $2 IS NOT NULL
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
RETURNING sell_req_id, sell_amount, currency_from, currency_to, tg_username, sell_by_card, sell_amount_by_card, sell_by_cash, sell_amount_by_cash, sell_exchange_rate, is_actual, created_at, updated_at, is_deleted, comment
`

type UpdateSellRequestParams struct {
	SellAmount       pgtype.Int8 `json:"sell_amount"`
	CurrencyFrom     pgtype.Text `json:"currency_from"`
	CurrencyTo       pgtype.Text `json:"currency_to"`
	TgUsername       pgtype.Text `json:"tg_username"`
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
		arg.SellAmount,
		arg.CurrencyFrom,
		arg.CurrencyTo,
		arg.TgUsername,
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
		&i.SellAmount,
		&i.CurrencyFrom,
		&i.CurrencyTo,
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
