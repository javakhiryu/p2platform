// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    telegram_id,
    tg_username,
    first_name,
    last_name
) VALUES (
    $1, $2, $3, $4
)
RETURNING telegram_id, tg_username, first_name, last_name, created_at
`

type CreateUserParams struct {
	TelegramID int64       `json:"telegram_id"`
	TgUsername string      `json:"tg_username"`
	FirstName  pgtype.Text `json:"first_name"`
	LastName   pgtype.Text `json:"last_name"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.TelegramID,
		arg.TgUsername,
		arg.FirstName,
		arg.LastName,
	)
	var i User
	err := row.Scan(
		&i.TelegramID,
		&i.TgUsername,
		&i.FirstName,
		&i.LastName,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE telegram_id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, telegramID int64) error {
	_, err := q.db.Exec(ctx, deleteUser, telegramID)
	return err
}

const getUser = `-- name: GetUser :one
SELECT telegram_id, tg_username, first_name, last_name, created_at FROM users WHERE telegram_id = $1
`

func (q *Queries) GetUser(ctx context.Context, telegramID int64) (User, error) {
	row := q.db.QueryRow(ctx, getUser, telegramID)
	var i User
	err := row.Scan(
		&i.TelegramID,
		&i.TgUsername,
		&i.FirstName,
		&i.LastName,
		&i.CreatedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET
    tg_username = COALESCE($1, tg_username),
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    updated_at = CASE
        WHEN $1 IS NOT NULL
          OR $2 IS NOT NULL
          OR $3 IS NOT NULL
        THEN now()
        ELSE updated_at
    END
WHERE telegram_id = $4
RETURNING telegram_id, tg_username, first_name, last_name, created_at
`

type UpdateUserParams struct {
	TgUsername pgtype.Text `json:"tg_username"`
	FirstName  pgtype.Text `json:"first_name"`
	LastName   pgtype.Text `json:"last_name"`
	TelegramID int64       `json:"telegram_id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.TgUsername,
		arg.FirstName,
		arg.LastName,
		arg.TelegramID,
	)
	var i User
	err := row.Scan(
		&i.TelegramID,
		&i.TgUsername,
		&i.FirstName,
		&i.LastName,
		&i.CreatedAt,
	)
	return i, err
}
