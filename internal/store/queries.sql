-- internal/store/queries.sql

-- name: GetUser :one
SELECT id, username, password_hash, created_at
FROM users 
WHERE username = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (username, password_hash)
VALUES ($1, $2)
RETURNING id, username, password_hash, created_at;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);