-- name: AddSpaceMember :one
INSERT INTO space_members (
    space_id,
    user_id,
    username
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetSpaceMembersCount :one
SELECT COUNT(*) FROM space_members WHERE space_id = $1;

-- name: IsUserInSameSpaceAsSeller :one
SELECT EXISTS (
    SELECT 1
    FROM space_members sm1
    JOIN space_members sm2 ON sm1.space_id = sm2.space_id
    WHERE sm1.user_id = $1 AND sm2.user_id = $2
);

-- name: GetSpaceMember :one
SELECT * FROM space_members WHERE space_id = $1 AND user_id = $2;

-- name: IsUserInSpace :one
SELECT EXISTS (
    SELECT 1 FROM space_members
    WHERE user_id = $1 AND space_id = $2
);

-- name: GetSpaceIdByUserId :many
SELECT space_id FROM space_members WHERE user_id = $1;

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