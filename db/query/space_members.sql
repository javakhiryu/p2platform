-- name: AddSpaceMember :one
INSERT INTO space_members (
    space_id,
    user_id,
    username
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetSpaceMember :one
SELECT * FROM space_members WHERE space_id = $1 AND user_id = $2;

-- name: ListSpaceMembersByUsernameAsc :many
SELECT *
FROM space_members
WHERE (username, user_id) > ($1, $2)
ORDER BY username ASC, user_id ASC
LIMIT $3;


-- name: ListSpaceMembersByUsernameDesc :many
SELECT *
FROM space_members
WHERE (username, user_id) < ($1, $2)
ORDER BY username DESC, user_id DESC
LIMIT $3;

-- name: DeleteSpaceMember :exec
DELETE FROM space_members WHERE space_id = $1 AND user_id = $2;