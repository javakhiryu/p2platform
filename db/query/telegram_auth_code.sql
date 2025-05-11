-- name: CreateTelegramAuthCode :one
INSERT INTO telegram_auth_codes (
    auth_code,
    expires_at
) VALUES (
    $1, $2
) 
RETURNING *;

-- name: GetTelegramAuthCode :one
SELECT * FROM telegram_auth_codes WHERE auth_code = $1;

-- name: ConfirmTelegramAuthCode :exec
UPDATE telegram_auth_codes
SET
    telegram_id = $2,
    status = $3
WHERE
    auth_code = $1
RETURNING *;

-- name: ListExpireAuthCodes :many
SELECT * FROM telegram_auth_codes
WHERE expires_at < now()
  AND status = 'pending';

-- name: ExpireTelegramAuthCode :exec
UPDATE telegram_auth_codes
SET
    status = 'expired'
WHERE
    auth_code = $1;
