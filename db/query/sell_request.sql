-- name: CreateSellRequest :one
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
RETURNING *;

-- name: GetSellRequestById :one
SELECT * FROM sell_requests WHERE sell_req_id = $1;

-- name: ListSellRequests :many
SELECT * FROM sell_requests
WHERE is_deleted = false
ORDER BY created_at ASC
LIMIT $1 
OFFSET $2;



-- name: CountOfSellRequests :one
SELECT COUNT(*) FROM sell_requests WHERE is_deleted = false;

-- name: ListSellRequestsByTelegramId :many
SELECT * FROM sell_requests
WHERE telegram_id = $1 AND is_deleted = false
ORDER BY created_at ASC
LIMIT $2 
OFFSET $3;

-- name: ListSellRequestsBySpaceAndTelegramID :many
SELECT sr.*
FROM sell_requests sr
JOIN space_members sm ON sr.telegram_id = sm.user_id
WHERE sr.telegram_id = $1
  AND sm.space_id = $2
  AND sr.is_deleted = false
ORDER BY sr.updated_at DESC
LIMIT $3
OFFSET $4;

-- name: ListSellRequestsBySpace :many
SELECT sr.*
FROM sell_requests sr
JOIN space_members sm ON sr.telegram_id = sm.user_id
WHERE sm.space_id = $1 AND sr.is_deleted = false
ORDER BY sr.updated_at DESC
LIMIT $2
OFFSET $3;


-- name: CountOfSellRequestsByTelegramId :one
SELECT COUNT(*) FROM sell_requests WHERE telegram_id = $1 AND is_deleted = false;

-- name: UpdateSellRequest :one
UPDATE sell_requests
SET
    sell_total_amount = COALESCE(sqlc.narg('sell_total_amount'), sell_total_amount),
    sell_money_source = COALESCE(sqlc.narg('sell_money_source'), sell_money_source),
    currency_from = COALESCE(sqlc.narg('currency_from'), currency_from),
    currency_to = COALESCE(sqlc.narg('currency_to'), currency_to),
    sell_by_card = COALESCE(sqlc.narg('sell_by_card'), sell_by_card),
    sell_amount_by_card = COALESCE(sqlc.narg('sell_amount_by_card'), sell_amount_by_card),
    sell_by_cash = COALESCE(sqlc.narg('sell_by_cash'), sell_by_cash),
    sell_amount_by_cash = COALESCE(sqlc.narg('sell_amount_by_cash'), sell_amount_by_cash),
    sell_exchange_rate = COALESCE(sqlc.narg('sell_exchange_rate'), sell_exchange_rate),
    comment = COALESCE(sqlc.narg('comment'), comment),
    updated_at = CASE
        WHEN sqlc.narg('sell_total_amount') IS NOT NULL
          OR sqlc.narg('currency_from') IS NOT NULL
          OR sqlc.narg('currency_to') IS NOT NULL
          OR sqlc.narg('sell_by_card') IS NOT NULL
          OR sqlc.narg('sell_amount_by_card') IS NOT NULL
          OR sqlc.narg('sell_by_cash') IS NOT NULL
          OR sqlc.narg('sell_amount_by_cash') IS NOT NULL
          OR sqlc.narg('sell_exchange_rate') IS NOT NULL
          OR sqlc.narg('comment') IS NOT NULL
        THEN now()
        ELSE updated_at
    END
WHERE sell_req_id = sqlc.arg('sell_req_id')
RETURNING *;


-- name: OpenCloseSellRequest :one
UPDATE sell_requests
SET
  is_actual = $1,
  updated_at = now()
WHERE
  sell_req_id = $2
RETURNING *;


-- name: DeleteSellRequest :one
UPDATE sell_requests
SET
  is_deleted = true,
  updated_at = now()
WHERE
  sell_req_id = $1
RETURNING is_deleted;

-- name: GetSellRequestForUpdate :one
SELECT * FROM sell_requests
WHERE sell_req_id = $1
FOR UPDATE;