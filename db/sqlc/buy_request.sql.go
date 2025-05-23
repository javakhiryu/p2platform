// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: buy_request.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const changeStateBuyRequest = `-- name: ChangeStateBuyRequest :one
UPDATE buy_requests
SET
  state = $1,
  state_updated_at = now()
WHERE
  buy_req_id = $2
RETURNING buy_req_id, sell_req_id, buy_total_amount, telegram_id, tg_username, buy_by_card, buy_amount_by_card, buy_by_cash, buy_amount_by_cash, close_confirm_by_seller, close_confirm_by_buyer, seller_confirmed_at, buyer_confirmed_at, state, state_updated_at, created_at, expires_at
`

type ChangeStateBuyRequestParams struct {
	State    string    `json:"state"`
	BuyReqID uuid.UUID `json:"buy_req_id"`
}

func (q *Queries) ChangeStateBuyRequest(ctx context.Context, arg ChangeStateBuyRequestParams) (BuyRequest, error) {
	row := q.db.QueryRow(ctx, changeStateBuyRequest, arg.State, arg.BuyReqID)
	var i BuyRequest
	err := row.Scan(
		&i.BuyReqID,
		&i.SellReqID,
		&i.BuyTotalAmount,
		&i.TelegramID,
		&i.TgUsername,
		&i.BuyByCard,
		&i.BuyAmountByCard,
		&i.BuyByCash,
		&i.BuyAmountByCash,
		&i.CloseConfirmBySeller,
		&i.CloseConfirmByBuyer,
		&i.SellerConfirmedAt,
		&i.BuyerConfirmedAt,
		&i.State,
		&i.StateUpdatedAt,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const closeBuyRequestBySellRequest = `-- name: CloseBuyRequestBySellRequest :exec
UPDATE buy_requests
SET
  state = 'closed',
  state_updated_at = now()
WHERE
  sell_req_id = $1 AND state = 'open'
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

const countBuyRequestsByUserInSpace = `-- name: CountBuyRequestsByUserInSpace :one
SELECT COUNT(*)
FROM buy_requests br
JOIN space_members sm ON br.telegram_id = sm.user_id AND br.space_id = sm.space_id
WHERE sm.user_id = $1 AND sm.space_id = $2 AND br.state = 'open'
`

type CountBuyRequestsByUserInSpaceParams struct {
	UserID  int64     `json:"user_id"`
	SpaceID uuid.UUID `json:"space_id"`
}

func (q *Queries) CountBuyRequestsByUserInSpace(ctx context.Context, arg CountBuyRequestsByUserInSpaceParams) (int64, error) {
	row := q.db.QueryRow(ctx, countBuyRequestsByUserInSpace, arg.UserID, arg.SpaceID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countOfBuyRequests = `-- name: CountOfBuyRequests :one
SELECT COUNT(*) FROM buy_requests WHERE sell_req_id = $1 AND state = 'open'
`

func (q *Queries) CountOfBuyRequests(ctx context.Context, sellReqID int32) (int64, error) {
	row := q.db.QueryRow(ctx, countOfBuyRequests, sellReqID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createBuyRequest = `-- name: CreateBuyRequest :one
INSERT INTO buy_requests (
  buy_req_id,
  sell_req_id,
  buy_total_amount,
  telegram_id,
  tg_username,
  buy_by_card,
  buy_amount_by_card,
  buy_by_cash,
  buy_amount_by_cash
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING buy_req_id, sell_req_id, buy_total_amount, telegram_id, tg_username, buy_by_card, buy_amount_by_card, buy_by_cash, buy_amount_by_cash, close_confirm_by_seller, close_confirm_by_buyer, seller_confirmed_at, buyer_confirmed_at, state, state_updated_at, created_at, expires_at
`

type CreateBuyRequestParams struct {
	BuyReqID        uuid.UUID   `json:"buy_req_id"`
	SellReqID       int32       `json:"sell_req_id"`
	BuyTotalAmount  int64       `json:"buy_total_amount"`
	TelegramID      int64       `json:"telegram_id"`
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
		arg.TelegramID,
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
		&i.TelegramID,
		&i.TgUsername,
		&i.BuyByCard,
		&i.BuyAmountByCard,
		&i.BuyByCash,
		&i.BuyAmountByCash,
		&i.CloseConfirmBySeller,
		&i.CloseConfirmByBuyer,
		&i.SellerConfirmedAt,
		&i.BuyerConfirmedAt,
		&i.State,
		&i.StateUpdatedAt,
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
SELECT buy_req_id, sell_req_id, buy_total_amount, telegram_id, tg_username, buy_by_card, buy_amount_by_card, buy_by_cash, buy_amount_by_cash, close_confirm_by_seller, close_confirm_by_buyer, seller_confirmed_at, buyer_confirmed_at, state, state_updated_at, created_at, expires_at FROM buy_requests WHERE buy_req_id = $1
`

func (q *Queries) GetBuyRequestById(ctx context.Context, buyReqID uuid.UUID) (BuyRequest, error) {
	row := q.db.QueryRow(ctx, getBuyRequestById, buyReqID)
	var i BuyRequest
	err := row.Scan(
		&i.BuyReqID,
		&i.SellReqID,
		&i.BuyTotalAmount,
		&i.TelegramID,
		&i.TgUsername,
		&i.BuyByCard,
		&i.BuyAmountByCard,
		&i.BuyByCash,
		&i.BuyAmountByCash,
		&i.CloseConfirmBySeller,
		&i.CloseConfirmByBuyer,
		&i.SellerConfirmedAt,
		&i.BuyerConfirmedAt,
		&i.State,
		&i.StateUpdatedAt,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const listBuyRequests = `-- name: ListBuyRequests :many
SELECT buy_req_id, sell_req_id, buy_total_amount, telegram_id, tg_username, buy_by_card, buy_amount_by_card, buy_by_cash, buy_amount_by_cash, close_confirm_by_seller, close_confirm_by_buyer, seller_confirmed_at, buyer_confirmed_at, state, state_updated_at, created_at, expires_at FROM buy_requests
WHERE sell_req_id = $1
ORDER BY created_at ASC
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
			&i.TelegramID,
			&i.TgUsername,
			&i.BuyByCard,
			&i.BuyAmountByCard,
			&i.BuyByCash,
			&i.BuyAmountByCash,
			&i.CloseConfirmBySeller,
			&i.CloseConfirmByBuyer,
			&i.SellerConfirmedAt,
			&i.BuyerConfirmedAt,
			&i.State,
			&i.StateUpdatedAt,
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

const listBuyRequestsByUserInSpace = `-- name: ListBuyRequestsByUserInSpace :many
SELECT br.buy_req_id, br.sell_req_id, br.buy_total_amount, br.telegram_id, br.tg_username, br.buy_by_card, br.buy_amount_by_card, br.buy_by_cash, br.buy_amount_by_cash, br.close_confirm_by_seller, br.close_confirm_by_buyer, br.seller_confirmed_at, br.buyer_confirmed_at, br.state, br.state_updated_at, br.created_at, br.expires_at
FROM buy_requests br
JOIN space_members sm ON br.telegram_id = sm.user_id AND br.space_id = sm.space_id
WHERE sm.user_id = $1 AND sm.space_id = $2
AND br.state = 'open'
ORDER BY br.created_at ASC
LIMIT $3 OFFSET $4
`

type ListBuyRequestsByUserInSpaceParams struct {
	UserID  int64     `json:"user_id"`
	SpaceID uuid.UUID `json:"space_id"`
	Limit   int32     `json:"limit"`
	Offset  int32     `json:"offset"`
}

func (q *Queries) ListBuyRequestsByUserInSpace(ctx context.Context, arg ListBuyRequestsByUserInSpaceParams) ([]BuyRequest, error) {
	rows, err := q.db.Query(ctx, listBuyRequestsByUserInSpace,
		arg.UserID,
		arg.SpaceID,
		arg.Limit,
		arg.Offset,
	)
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
			&i.TelegramID,
			&i.TgUsername,
			&i.BuyByCard,
			&i.BuyAmountByCard,
			&i.BuyByCash,
			&i.BuyAmountByCash,
			&i.CloseConfirmBySeller,
			&i.CloseConfirmByBuyer,
			&i.SellerConfirmedAt,
			&i.BuyerConfirmedAt,
			&i.State,
			&i.StateUpdatedAt,
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

const listExpiredBuyRequests = `-- name: ListExpiredBuyRequests :many
SELECT buy_req_id, sell_req_id, buy_total_amount, telegram_id, tg_username, buy_by_card, buy_amount_by_card, buy_by_cash, buy_amount_by_cash, close_confirm_by_seller, close_confirm_by_buyer, seller_confirmed_at, buyer_confirmed_at, state, state_updated_at, created_at, expires_at FROM buy_requests
WHERE expires_at < now()
  AND state = 'open'
`

func (q *Queries) ListExpiredBuyRequests(ctx context.Context) ([]BuyRequest, error) {
	rows, err := q.db.Query(ctx, listExpiredBuyRequests)
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
			&i.TelegramID,
			&i.TgUsername,
			&i.BuyByCard,
			&i.BuyAmountByCard,
			&i.BuyByCash,
			&i.BuyAmountByCash,
			&i.CloseConfirmBySeller,
			&i.CloseConfirmByBuyer,
			&i.SellerConfirmedAt,
			&i.BuyerConfirmedAt,
			&i.State,
			&i.StateUpdatedAt,
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
