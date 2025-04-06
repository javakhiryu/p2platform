-- name: CreateBuyRequest :one
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
RETURNING *;

-- name: GetBuyRequestById :one
SELECT * FROM buy_requests WHERE buy_req_id = $1;

-- name: ListBuyRequests :many
SELECT * FROM buy_requests
WHERE sell_req_id = $1
    AND is_successful = false
ORDER BY created_at DESC
LIMIT $2 
OFFSET $3;

-- name: UpdateBuyRequest :one
UPDATE buy_requests
SET tg_username= $1
WHERE buy_req_id = $2
RETURNING *;


-- name: OpenCloseBuyRequest :one
UPDATE buy_requests
SET
  is_successful = $1
WHERE
  buy_req_id = $2
RETURNING *;

-- name: DeleteBuyRequest :exec
DELETE FROM buy_requests
WHERE
  buy_req_id = $1;

