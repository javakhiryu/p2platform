-- name: CreateUser :one
INSERT INTO users (
    telegram_id,
    tg_username,
    first_name,
    last_name
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE telegram_id = $1;

-- name: UpdateUser :one
UPDATE users
SET
    tg_username = COALESCE(sqlc.narg('tg_username'), tg_username),
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    updated_at = CASE
        WHEN sqlc.narg('tg_username') IS NOT NULL
          OR sqlc.narg('first_name') IS NOT NULL
          OR sqlc.narg('last_name') IS NOT NULL
        THEN now()
        ELSE updated_at
    END
WHERE telegram_id = sqlc.arg('telegram_id')
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE telegram_id = $1;