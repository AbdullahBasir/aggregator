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
SELECT users.name, feeds.url, feeds.name
FROM users
INNER JOIN feeds ON users.id = feeds.user_id;

-- name: ListUsers :many
SELECT * FROM users;

-- name: ResetUser :exec
DELETE FROM users;