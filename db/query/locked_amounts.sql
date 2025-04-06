-- name: CreateLockedAmount :one
INSERT INTO locked_amounts (
  sell_req_id,
  buy_req_id,
  locked_total,
  locked_by_card,
  locked_by_cash
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetLockedAmountByBuyReqID :one
SELECT * FROM locked_amounts
WHERE buy_req_id = $1;

-- name: ListLockedAmountsBySellReqID :many
SELECT * FROM locked_amounts
WHERE sell_req_id = $1
ORDER BY created_at;

-- name: ReleaseLockedAmountByBuyRequest :one
UPDATE locked_amounts
SET
  is_released = true,
  released_at = now()
WHERE
  buy_req_id = $1
RETURNING *;

-- name: DeleteLockedAmountByBuyReqID :exec
DELETE FROM locked_amounts
WHERE buy_req_id = $1;

-- name: ReleaseLockedAmountsBySellRequest :exec
UPDATE locked_amounts
SET is_released = true,
    released_at = now()
WHERE sell_req_id = $1 AND is_released = false;

