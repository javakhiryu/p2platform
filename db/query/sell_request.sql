-- name: CreateSellRequest :one
INSERT INTO sell_requests (
  sell_amount,
  currency,
  tg_username,
  sell_by_card,
  sell_amount_by_card,
  sell_by_cash,
  sell_amunt_by_cash,
  sell_exchange_rate,
  comment
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetSellRequestById :one
SELECT * FROM sell_requests WHERE sell_req_id = $1;

-- name: ListSellRequests :many
SELECT * FROM sell_requests
WHERE is_deleted = false
ORDER BY created_at DESC
LIMIT $1 
OFFSET $2;

-- name: UpdateSellRequest :one
UPDATE sell_requests
SET
  sell_amount = coalesce($1, sell_amount),
    currency = coalesce($2, currency),
    tg_username = coalesce($3, tg_username),
    sell_by_card = coalesce($4, sell_by_card),
    sell_amount_by_card = coalesce($5, sell_amount_by_card),
    sell_by_cash = coalesce($6, sell_by_cash),
    sell_amunt_by_cash = coalesce($7, sell_amunt_by_cash),
    sell_exchange_rate = coalesce($8, sell_exchange_rate),
    comment = coalesce($9, comment),
    is_actual = coalesce($10, is_actual),
    -- set updated_at to now() if any of the fields are updated
    -- otherwise keep the old value
    -- this is a workaround for the fact that we can't use
    -- coalesce on the updated_at field
    updated_at = CASE
      WHEN $1 IS NOT NULL OR $2 IS NOT NULL OR $3 IS NOT NULL OR $4 IS NOT NULL OR $5 IS NOT NULL OR $6 IS NOT NULL OR $7 IS NOT NULL OR $8 IS NOT NULL OR $9 IS NOT NULL OR $10 IS NOT NULL OR $11 IS NOT NULL
      THEN now()
      ELSE updated_at
    END
WHERE sell_req_id = $11
RETURNING *;

-- name: DeleteSellRequest :one
UPDATE sell_requests
SET
  is_deleted = true,
  updated_at = now()
WHERE
  sell_req_id = $1
RETURNING *;