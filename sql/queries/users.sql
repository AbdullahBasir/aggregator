-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users 
WHERE name = $1;

-- name: GetFeeds :many
SELECT users.name, feed.url, feed.name
FROM users
INNER JOIN feed ON users.id = feed.user_id;

-- name: ListUsers :many
SELECT * FROM users;

-- name: ResetUser :exec
DELETE FROM users;