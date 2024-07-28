-- name: CreateVerifyAdminEmail :one
INSERT INTO admin_verify_emails (username, email, secret_code)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetVerifyAdminEmail :one
SELECT id, username, email, is_used, secret_code, created_at
FROM admin_verify_emails
WHERE email = $1 and secret_code = $2;

-- name: UpdateVerifyAdminEmail :exec
UPDATE admin_verify_emails
SET is_used = COALESCE(sqlc.arg(is_used), is_used)
WHERE email = $1 AND secret_code = $2
RETURNING *;