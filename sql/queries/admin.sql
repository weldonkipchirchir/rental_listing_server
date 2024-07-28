-- name: CreateAdmin :one
INSERT INTO admins (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id, username, email, created_at;

-- name: GetAdminByID :one
SELECT id, username, email, created_at
FROM admins
WHERE id = $1;

-- name: UpdateAdminPassword :exec
UPDATE admins
SET
    password_hash = $2,
    username = $3
WHERE id = $1;

-- name: UpdateAdminEmailVerified :exec
UPDATE admins
SET is_email_verified = TRUE
WHERE email = $1;

-- name: UpdateAdminForgotPassword :exec
UPDATE admins
SET password_hash = $2
WHERE email = $1;

-- name: DeleteAdmin :exec
DELETE FROM admins
WHERE id = $1;

-- name: GetAdmin :one
SELECT * FROM admins
WHERE email = $1 LIMIT 1;