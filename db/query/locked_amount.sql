-- name: CreateLockedAmount :one
INSERT INTO locked_amounts (
  sell_req_id,
  buy_req_id,
  locked_total_amount,
  locked_by_card,
  locked_by_cash
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetLockedAmount :one
SELECT * FROM locked_amounts
WHERE buy_req_id = $1;

-- name: GetLockedAmountBySellRequest :many
SELECT * FROM locked_amounts
WHERE sell_req_id = $1
AND is_released = false;

-- name: ListLockedAmounts :many
SELECT * FROM locked_amounts
WHERE sell_req_id = $1
ORDER BY created_at
LIMIT $2
OFFSET $3;

-- name: ReleaseLockedAmountByBuyRequest :exec
UPDATE locked_amounts
SET
  is_released = true,
  released_at = now()
WHERE
  buy_req_id = $1
  AND is_released = false;

-- name: ReleaseLockedAmountsBySellRequest :exec
UPDATE locked_amounts
SET is_released = true,
    released_at = now()
WHERE sell_req_id = $1 AND is_released = false;

