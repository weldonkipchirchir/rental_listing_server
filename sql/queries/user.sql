-- name: CreateUser :one
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id, username, email, created_at;

-- name: GetUserByID :one
SELECT id, username, email, created_at
FROM users
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET
    password_hash = $2,
    username = $3
WHERE id = $1;

-- name: UpdateUserEmailVerified :exec
UPDATE users
SET is_email_verified = TRUE
WHERE email = $1;

-- name: UpdateUserForgotPassword :exec
UPDATE users
SET password_hash = $2
WHERE email = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;