// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: space_members.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const addSpaceMember = `-- name: AddSpaceMember :one
INSERT INTO space_members (
    space_id,
    user_id,
    username
) VALUES (
    $1, $2, $3
)
RETURNING space_id, user_id, username, joined_at
`

type AddSpaceMemberParams struct {
	SpaceID  uuid.UUID `json:"space_id"`
	UserID   int64     `json:"user_id"`
	Username string    `json:"username"`
}

func (q *Queries) AddSpaceMember(ctx context.Context, arg AddSpaceMemberParams) (SpaceMember, error) {
	row := q.db.QueryRow(ctx, addSpaceMember, arg.SpaceID, arg.UserID, arg.Username)
	var i SpaceMember
	err := row.Scan(
		&i.SpaceID,
		&i.UserID,
		&i.Username,
		&i.JoinedAt,
	)
	return i, err
}

const deleteSpaceMember = `-- name: DeleteSpaceMember :exec
DELETE FROM space_members WHERE space_id = $1 AND user_id = $2
`

type DeleteSpaceMemberParams struct {
	SpaceID uuid.UUID `json:"space_id"`
	UserID  int64     `json:"user_id"`
}

func (q *Queries) DeleteSpaceMember(ctx context.Context, arg DeleteSpaceMemberParams) error {
	_, err := q.db.Exec(ctx, deleteSpaceMember, arg.SpaceID, arg.UserID)
	return err
}

const getSpaceIdByUserId = `-- name: GetSpaceIdByUserId :many
SELECT space_id FROM space_members WHERE user_id = $1
`

func (q *Queries) GetSpaceIdByUserId(ctx context.Context, userID int64) ([]uuid.UUID, error) {
	rows, err := q.db.Query(ctx, getSpaceIdByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []uuid.UUID{}
	for rows.Next() {
		var space_id uuid.UUID
		if err := rows.Scan(&space_id); err != nil {
			return nil, err
		}
		items = append(items, space_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSpaceMember = `-- name: GetSpaceMember :one
SELECT space_id, user_id, username, joined_at FROM space_members WHERE space_id = $1 AND user_id = $2
`

type GetSpaceMemberParams struct {
	SpaceID uuid.UUID `json:"space_id"`
	UserID  int64     `json:"user_id"`
}

func (q *Queries) GetSpaceMember(ctx context.Context, arg GetSpaceMemberParams) (SpaceMember, error) {
	row := q.db.QueryRow(ctx, getSpaceMember, arg.SpaceID, arg.UserID)
	var i SpaceMember
	err := row.Scan(
		&i.SpaceID,
		&i.UserID,
		&i.Username,
		&i.JoinedAt,
	)
	return i, err
}

const getSpaceMembersCount = `-- name: GetSpaceMembersCount :one
SELECT COUNT(*) FROM space_members WHERE space_id = $1
`

func (q *Queries) GetSpaceMembersCount(ctx context.Context, spaceID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, getSpaceMembersCount, spaceID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const isUserInSameSpaceAsSeller = `-- name: IsUserInSameSpaceAsSeller :one
SELECT EXISTS (
    SELECT 1
    FROM space_members sm1
    JOIN space_members sm2 ON sm1.space_id = sm2.space_id
    WHERE sm1.user_id = $1 AND sm2.user_id = $2
)
`

type IsUserInSameSpaceAsSellerParams struct {
	UserID   int64 `json:"user_id"`
	UserID_2 int64 `json:"user_id_2"`
}

func (q *Queries) IsUserInSameSpaceAsSeller(ctx context.Context, arg IsUserInSameSpaceAsSellerParams) (bool, error) {
	row := q.db.QueryRow(ctx, isUserInSameSpaceAsSeller, arg.UserID, arg.UserID_2)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const isUserInSpace = `-- name: IsUserInSpace :one
SELECT EXISTS (
    SELECT 1 FROM space_members
    WHERE user_id = $1 AND space_id = $2
)
`

type IsUserInSpaceParams struct {
	UserID  int64     `json:"user_id"`
	SpaceID uuid.UUID `json:"space_id"`
}

func (q *Queries) IsUserInSpace(ctx context.Context, arg IsUserInSpaceParams) (bool, error) {
	row := q.db.QueryRow(ctx, isUserInSpace, arg.UserID, arg.SpaceID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const listSpaceMembersByUsernameAsc = `-- name: ListSpaceMembersByUsernameAsc :many
SELECT space_id, user_id, username, joined_at
FROM space_members
WHERE (username, user_id) > ($1, $2)
ORDER BY username ASC, user_id ASC
LIMIT $3
`

type ListSpaceMembersByUsernameAscParams struct {
	Username   string `json:"username"`
	Username_2 string `json:"username_2"`
	Limit      int32  `json:"limit"`
}

func (q *Queries) ListSpaceMembersByUsernameAsc(ctx context.Context, arg ListSpaceMembersByUsernameAscParams) ([]SpaceMember, error) {
	rows, err := q.db.Query(ctx, listSpaceMembersByUsernameAsc, arg.Username, arg.Username_2, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []SpaceMember{}
	for rows.Next() {
		var i SpaceMember
		if err := rows.Scan(
			&i.SpaceID,
			&i.UserID,
			&i.Username,
			&i.JoinedAt,
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

const listSpaceMembersByUsernameDesc = `-- name: ListSpaceMembersByUsernameDesc :many
SELECT space_id, user_id, username, joined_at
FROM space_members
WHERE (username, user_id) < ($1, $2)
ORDER BY username DESC, user_id DESC
LIMIT $3
`

type ListSpaceMembersByUsernameDescParams struct {
	Username   string `json:"username"`
	Username_2 string `json:"username_2"`
	Limit      int32  `json:"limit"`
}

func (q *Queries) ListSpaceMembersByUsernameDesc(ctx context.Context, arg ListSpaceMembersByUsernameDescParams) ([]SpaceMember, error) {
	rows, err := q.db.Query(ctx, listSpaceMembersByUsernameDesc, arg.Username, arg.Username_2, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []SpaceMember{}
	for rows.Next() {
		var i SpaceMember
		if err := rows.Scan(
			&i.SpaceID,
			&i.UserID,
			&i.Username,
			&i.JoinedAt,
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
