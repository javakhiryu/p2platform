-- name: CreateBuyRequest :one
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
RETURNING *;

-- name: GetBuyRequestById :one
SELECT * FROM buy_requests WHERE buy_req_id = $1;

-- name: ListBuyRequests :many
SELECT * FROM buy_requests
WHERE sell_req_id = $1
ORDER BY created_at ASC
LIMIT $2 
OFFSET $3;

-- name: ListBuyRequestsByTelegramId :many
SELECT * FROM buy_requests
WHERE telegram_id = $1
ORDER BY created_at ASC
LIMIT $2 
OFFSET $3;

-- name: CloseConfirmBySeller :exec
UPDATE buy_requests
SET
  close_confirm_by_seller = $1,
  seller_confirmed_at = now()
WHERE
  buy_req_id = $2;

-- name: CloseConfirmByBuyer :exec
UPDATE buy_requests
SET
  close_confirm_by_buyer = $1,
  buyer_confirmed_at = now()
WHERE
  buy_req_id = $2;

-- name: OpenCloseBuyRequest :one
UPDATE buy_requests
SET
  is_closed = $1,
  closed_at = now()
WHERE
  buy_req_id = $2
RETURNING *;

-- name: CloseBuyRequestBySellRequest :exec
UPDATE buy_requests
SET
  is_closed = true,
  closed_at = now()
WHERE
  sell_req_id = $1 AND is_closed = false;

-- name: DeleteBuyRequest :exec
DELETE FROM buy_requests WHERE buy_req_id = $1;

-- name: ListExpiredBuyRequests :many
SELECT * FROM buy_requests
WHERE expires_at < now()
  AND is_closed = false;

