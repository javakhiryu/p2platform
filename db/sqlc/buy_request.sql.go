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

const closeBuyRequestBySellRequest = `-- name: CloseBuyRequestBySellRequest :exec
UPDATE buy_requests
SET
  is_closed = true,
  closed_at = now()
WHERE
  sell_req_id = $1 AND is_closed = false
`

func (q *Queries) CloseBuyRequestBySellRequest(ctx context.Context, sellReqID int32) error {
	_, err := q.db.Exec(ctx, closeBuyRequestBySellRequest, sellReqID)
	return err
}

const closeConfirmByBuyer = `-- name: CloseConfirmByBuyer :exec
UPDATE buy_requests
SET
  close_confirm_by_buyer = $1,
  buyer_confirmed_at = now()
WHERE
  buy_req_id = $2
`

type CloseConfirmByBuyerParams struct {
	CloseConfirmByBuyer pgtype.Bool `json:"close_confirm_by_buyer"`
	BuyReqID            uuid.UUID   `json:"buy_req_id"`
}

func (q *Queries) CloseConfirmByBuyer(ctx context.Context, arg CloseConfirmByBuyerParams) error {
	_, err := q.db.Exec(ctx, closeConfirmByBuyer, arg.CloseConfirmByBuyer, arg.BuyReqID)
	return err
}

const closeConfirmBySeller = `-- name: CloseConfirmBySeller :exec
UPDATE buy_requests
SET
  close_confirm_by_seller = $1,
  seller_confirmed_at = now()
WHERE
  buy_req_id = $2
`

type CloseConfirmBySellerParams struct {
	CloseConfirmBySeller pgtype.Bool `json:"close_confirm_by_seller"`
	BuyReqID             uuid.UUID   `json:"buy_req_id"`
}

func (q *Queries) CloseConfirmBySeller(ctx context.Context, arg CloseConfirmBySellerParams) error {
	_, err := q.db.Exec(ctx, closeConfirmBySeller, arg.CloseConfirmBySeller, arg.BuyReqID)
	return err
}

const createBuyRequest = `-- name: CreateBuyRequest :one
INSERT INTO buy_requests (
  buy_req_id,
  sell_req_id,
  buy_total_amount,
  tg_username,
  buy_by_card,
  buy_amount_by_card,
  buy_by_cash,
  buy_amount_by_cash
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING buy_req_id, sell_req_id, buy_total_amount, tg_username, buy_by_card, buy_amount_by_card, buy_by_cash, buy_amount_by_cash, close_confirm_by_seller, close_confirm_by_buyer, seller_confirmed_at, buyer_confirmed_at, is_closed, closed_at, created_at, expires_at
`

type CreateBuyRequestParams struct {
	BuyReqID        uuid.UUID   `json:"buy_req_id"`
	SellReqID       int32       `json:"sell_req_id"`
	BuyTotalAmount  int64       `json:"buy_total_amount"`
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
		arg.BuyTotalAmount,
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
		&i.BuyTotalAmount,
		&i.TgUsername,
		&i.BuyByCard,
		&i.BuyAmountByCard,
		&i.BuyByCash,
		&i.BuyAmountByCash,
		&i.CloseConfirmBySeller,
		&i.CloseConfirmByBuyer,
		&i.SellerConfirmedAt,
		&i.BuyerConfirmedAt,
		&i.IsClosed,
		&i.ClosedAt,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const deleteBuyRequest = `-- name: DeleteBuyRequest :exec
DELETE FROM buy_requests WHERE buy_req_id = $1
`

func (q *Queries) DeleteBuyRequest(ctx context.Context, buyReqID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteBuyRequest, buyReqID)
	return err
}

const getBuyRequestById = `-- name: GetBuyRequestById :one
SELECT buy_req_id, sell_req_id, buy_total_amount, tg_username, buy_by_card, buy_amount_by_card, buy_by_cash, buy_amount_by_cash, close_confirm_by_seller, close_confirm_by_buyer, seller_confirmed_at, buyer_confirmed_at, is_closed, closed_at, created_at, expires_at FROM buy_requests WHERE buy_req_id = $1
`

func (q *Queries) GetBuyRequestById(ctx context.Context, buyReqID uuid.UUID) (BuyRequest, error) {
	row := q.db.QueryRow(ctx, getBuyRequestById, buyReqID)
	var i BuyRequest
	err := row.Scan(
		&i.BuyReqID,
		&i.SellReqID,
		&i.BuyTotalAmount,
		&i.TgUsername,
		&i.BuyByCard,
		&i.BuyAmountByCard,
		&i.BuyByCash,
		&i.BuyAmountByCash,
		&i.CloseConfirmBySeller,
		&i.CloseConfirmByBuyer,
		&i.SellerConfirmedAt,
		&i.BuyerConfirmedAt,
		&i.IsClosed,
		&i.ClosedAt,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const listBuyRequests = `-- name: ListBuyRequests :many
SELECT buy_req_id, sell_req_id, buy_total_amount, tg_username, buy_by_card, buy_amount_by_card, buy_by_cash, buy_amount_by_cash, close_confirm_by_seller, close_confirm_by_buyer, seller_confirmed_at, buyer_confirmed_at, is_closed, closed_at, created_at, expires_at FROM buy_requests
WHERE sell_req_id = $1
    AND is_closed = false
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
			&i.BuyTotalAmount,
			&i.TgUsername,
			&i.BuyByCard,
			&i.BuyAmountByCard,
			&i.BuyByCash,
			&i.BuyAmountByCash,
			&i.CloseConfirmBySeller,
			&i.CloseConfirmByBuyer,
			&i.SellerConfirmedAt,
			&i.BuyerConfirmedAt,
			&i.IsClosed,
			&i.ClosedAt,
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
  is_closed = $1,
  closed_at = now()
WHERE
  buy_req_id = $2
  AND close_confirm_by_buyer = true
  AND close_confirm_by_seller = true
RETURNING buy_req_id, sell_req_id, buy_total_amount, tg_username, buy_by_card, buy_amount_by_card, buy_by_cash, buy_amount_by_cash, close_confirm_by_seller, close_confirm_by_buyer, seller_confirmed_at, buyer_confirmed_at, is_closed, closed_at, created_at, expires_at
`

type OpenCloseBuyRequestParams struct {
	IsClosed pgtype.Bool `json:"is_closed"`
	BuyReqID uuid.UUID   `json:"buy_req_id"`
}

func (q *Queries) OpenCloseBuyRequest(ctx context.Context, arg OpenCloseBuyRequestParams) (BuyRequest, error) {
	row := q.db.QueryRow(ctx, openCloseBuyRequest, arg.IsClosed, arg.BuyReqID)
	var i BuyRequest
	err := row.Scan(
		&i.BuyReqID,
		&i.SellReqID,
		&i.BuyTotalAmount,
		&i.TgUsername,
		&i.BuyByCard,
		&i.BuyAmountByCard,
		&i.BuyByCash,
		&i.BuyAmountByCash,
		&i.CloseConfirmBySeller,
		&i.CloseConfirmByBuyer,
		&i.SellerConfirmedAt,
		&i.BuyerConfirmedAt,
		&i.IsClosed,
		&i.ClosedAt,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const updateBuyRequest = `-- name: UpdateBuyRequest :one
UPDATE buy_requests
SET tg_username= $1
WHERE buy_req_id = $2
RETURNING buy_req_id, sell_req_id, buy_total_amount, tg_username, buy_by_card, buy_amount_by_card, buy_by_cash, buy_amount_by_cash, close_confirm_by_seller, close_confirm_by_buyer, seller_confirmed_at, buyer_confirmed_at, is_closed, closed_at, created_at, expires_at
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
		&i.BuyTotalAmount,
		&i.TgUsername,
		&i.BuyByCard,
		&i.BuyAmountByCard,
		&i.BuyByCash,
		&i.BuyAmountByCash,
		&i.CloseConfirmBySeller,
		&i.CloseConfirmByBuyer,
		&i.SellerConfirmedAt,
		&i.BuyerConfirmedAt,
		&i.IsClosed,
		&i.ClosedAt,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}
