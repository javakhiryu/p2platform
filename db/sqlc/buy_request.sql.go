// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: buy_request.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createBuyRequest = `-- name: CreateBuyRequest :one
INSERT INTO buy_requests (
   buy_req_id,
  sell_req_id,
  buy_amount,
  tg_username,
  buy_by_card,
  buy_amount_by_card,
  buy_by_cash,
  buy_amount_by_cash
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING buy_req_id, sell_req_id, buy_amount, tg_username, buy_by_card, buy_amount_by_card, buy_by_cash, buy_amount_by_cash, is_successful, created_at, expires_at
`

type CreateBuyRequestParams struct {
	BuyReqID        uuid.UUID   `json:"buy_req_id"`
	SellReqID       int32       `json:"sell_req_id"`
	BuyAmount       int64       `json:"buy_amount"`
	TgUsername      string      `json:"tg_username"`
	BuyByCard       pgtype.Bool `json:"buy_by_card"`
	BuyAmountByCard pgtype.Int8 `json:"buy_amount_by_card"`
	BuyByCash       pgtype.Bool `json:"buy_by_cash"`
	BuyAmountByCash pgtype.Int8 `json:"buy_amount_by_cash"`
}

func (q *Queries) CreateBuyRequest(ctx context.Context, arg CreateBuyRequestParams) (BuyRequest, error) {
	row := q.db.QueryRow(ctx, createBuyRequest,
		arg.BuyReqID,
		arg.SellReqID,
		arg.BuyAmount,
		arg.TgUsername,
		arg.BuyByCard,
		arg.BuyAmountByCard,
		arg.BuyByCash,
		arg.BuyAmountByCash,
	)
	var i BuyRequest
	err := row.Scan(
		&i.BuyReqID,
		&i.SellReqID,
		&i.BuyAmount,
		&i.TgUsername,
		&i.BuyByCard,
		&i.BuyAmountByCard,
		&i.BuyByCash,
		&i.BuyAmountByCash,
		&i.IsSuccessful,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const deleteBuyRequest = `-- name: DeleteBuyRequest :exec
DELETE FROM buy_requests
WHERE
  buy_req_id = $1
`

func (q *Queries) DeleteBuyRequest(ctx context.Context, buyReqID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteBuyRequest, buyReqID)
	return err
}

const getBuyRequestById = `-- name: GetBuyRequestById :one
SELECT buy_req_id, sell_req_id, buy_amount, tg_username, buy_by_card, buy_amount_by_card, buy_by_cash, buy_amount_by_cash, is_successful, created_at, expires_at FROM buy_requests WHERE buy_req_id = $1
`

func (q *Queries) GetBuyRequestById(ctx context.Context, buyReqID uuid.UUID) (BuyRequest, error) {
	row := q.db.QueryRow(ctx, getBuyRequestById, buyReqID)
	var i BuyRequest
	err := row.Scan(
		&i.BuyReqID,
		&i.SellReqID,
		&i.BuyAmount,
		&i.TgUsername,
		&i.BuyByCard,
		&i.BuyAmountByCard,
		&i.BuyByCash,
		&i.BuyAmountByCash,
		&i.IsSuccessful,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const listBuyRequests = `-- name: ListBuyRequests :many
SELECT buy_req_id, sell_req_id, buy_amount, tg_username, buy_by_card, buy_amount_by_card, buy_by_cash, buy_amount_by_cash, is_successful, created_at, expires_at FROM buy_requests
WHERE sell_req_id = $1
    AND is_successful = false
ORDER BY created_at DESC
LIMIT $2 
OFFSET $3
`

type ListBuyRequestsParams struct {
	SellReqID int32 `json:"sell_req_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListBuyRequests(ctx context.Context, arg ListBuyRequestsParams) ([]BuyRequest, error) {
	rows, err := q.db.Query(ctx, listBuyRequests, arg.SellReqID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []BuyRequest{}
	for rows.Next() {
		var i BuyRequest
		if err := rows.Scan(
			&i.BuyReqID,
			&i.SellReqID,
			&i.BuyAmount,
			&i.TgUsername,
			&i.BuyByCard,
			&i.BuyAmountByCard,
			&i.BuyByCash,
			&i.BuyAmountByCash,
			&i.IsSuccessful,
			&i.CreatedAt,
			&i.ExpiresAt,
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

const openCloseBuyRequest = `-- name: OpenCloseBuyRequest :one
UPDATE buy_requests
SET
  is_successful = $1
WHERE
  buy_req_id = $2
RETURNING buy_req_id, sell_req_id, buy_amount, tg_username, buy_by_card, buy_amount_by_card, buy_by_cash, buy_amount_by_cash, is_successful, created_at, expires_at
`

type OpenCloseBuyRequestParams struct {
	IsSuccessful pgtype.Bool `json:"is_successful"`
	BuyReqID     uuid.UUID   `json:"buy_req_id"`
}

func (q *Queries) OpenCloseBuyRequest(ctx context.Context, arg OpenCloseBuyRequestParams) (BuyRequest, error) {
	row := q.db.QueryRow(ctx, openCloseBuyRequest, arg.IsSuccessful, arg.BuyReqID)
	var i BuyRequest
	err := row.Scan(
		&i.BuyReqID,
		&i.SellReqID,
		&i.BuyAmount,
		&i.TgUsername,
		&i.BuyByCard,
		&i.BuyAmountByCard,
		&i.BuyByCash,
		&i.BuyAmountByCash,
		&i.IsSuccessful,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const updateBuyRequest = `-- name: UpdateBuyRequest :one
UPDATE buy_requests
SET tg_username= $1
WHERE buy_req_id = $2
RETURNING buy_req_id, sell_req_id, buy_amount, tg_username, buy_by_card, buy_amount_by_card, buy_by_cash, buy_amount_by_cash, is_successful, created_at, expires_at
`

type UpdateBuyRequestParams struct {
	TgUsername string    `json:"tg_username"`
	BuyReqID   uuid.UUID `json:"buy_req_id"`
}

func (q *Queries) UpdateBuyRequest(ctx context.Context, arg UpdateBuyRequestParams) (BuyRequest, error) {
	row := q.db.QueryRow(ctx, updateBuyRequest, arg.TgUsername, arg.BuyReqID)
	var i BuyRequest
	err := row.Scan(
		&i.BuyReqID,
		&i.SellReqID,
		&i.BuyAmount,
		&i.TgUsername,
		&i.BuyByCard,
		&i.BuyAmountByCard,
		&i.BuyByCash,
		&i.BuyAmountByCash,
		&i.IsSuccessful,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}
