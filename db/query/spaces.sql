-- name: CreateSpace :one
INSERT INTO spaces (
    space_id,
    space_name,
    hashed_password,
    description,
    creator_id
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetSpaceBySpaceId :one
SELECT * FROM spaces WHERE space_id = $1;

-- name: GetSpaceByCreatorId :many
SELECT * FROM spaces WHERE creator_id = $1;

-- name: ListSpacesAfterCursorByNameAsc :many
SELECT *
FROM spaces
WHERE (space_name, space_id) > ($1, $2) -- $1 = lastSpaceName (VARCHAR), $2 = lastSpaceId (UUID)
ORDER BY space_name ASC, space_id ASC
LIMIT $3;

-- name: ListSpacesAfterCursorByNameDesc :many
SELECT *
FROM spaces
WHERE (space_name, space_id) < ($1, $2) 
ORDER BY space_name DESC, space_id DESC
LIMIT $3;

-- name: UpdateSpaceInfo :one
UPDATE spaces
SET
    space_name = COALESCE(sqlc.narg('space_name'), space_name),
    description = COALESCE(sqlc.narg('description'), description),
    updated_at = CASE
        WHEN sqlc.narg('space_name') IS NOT NULL
          OR sqlc.narg('description') IS NOT NULL
        THEN now()
        ELSE updated_at
    END
WHERE space_id = sqlc.arg('space_id')
RETURNING *;

-- name: DeleteSpace :exec
DELETE FROM spaces WHERE space_id = $1;
