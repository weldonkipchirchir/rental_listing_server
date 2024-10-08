// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: verify_admin_email.sql

package db

import (
	"context"
	"time"
)

const createVerifyAdminEmail = `-- name: CreateVerifyAdminEmail :one
INSERT INTO admin_verify_emails (username, email, secret_code)
VALUES ($1, $2, $3)
RETURNING id, username, email, secret_code, is_used, created_at, expired_at
`

type CreateVerifyAdminEmailParams struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	SecretCode string `json:"secret_code"`
}

func (q *Queries) CreateVerifyAdminEmail(ctx context.Context, arg CreateVerifyAdminEmailParams) (AdminVerifyEmail, error) {
	row := q.db.QueryRowContext(ctx, createVerifyAdminEmail, arg.Username, arg.Email, arg.SecretCode)
	var i AdminVerifyEmail
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.SecretCode,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}

const getVerifyAdminEmail = `-- name: GetVerifyAdminEmail :one
SELECT id, username, email, is_used, secret_code, created_at
FROM admin_verify_emails
WHERE email = $1 and secret_code = $2
`

type GetVerifyAdminEmailParams struct {
	Email      string `json:"email"`
	SecretCode string `json:"secret_code"`
}

type GetVerifyAdminEmailRow struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	IsUsed     bool      `json:"is_used"`
	SecretCode string    `json:"secret_code"`
	CreatedAt  time.Time `json:"created_at"`
}

func (q *Queries) GetVerifyAdminEmail(ctx context.Context, arg GetVerifyAdminEmailParams) (GetVerifyAdminEmailRow, error) {
	row := q.db.QueryRowContext(ctx, getVerifyAdminEmail, arg.Email, arg.SecretCode)
	var i GetVerifyAdminEmailRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.IsUsed,
		&i.SecretCode,
		&i.CreatedAt,
	)
	return i, err
}

const updateVerifyAdminEmail = `-- name: UpdateVerifyAdminEmail :exec
UPDATE admin_verify_emails
SET is_used = COALESCE($3, is_used)
WHERE email = $1 AND secret_code = $2
RETURNING id, username, email, secret_code, is_used, created_at, expired_at
`

type UpdateVerifyAdminEmailParams struct {
	Email      string `json:"email"`
	SecretCode string `json:"secret_code"`
	IsUsed     bool   `json:"is_used"`
}

func (q *Queries) UpdateVerifyAdminEmail(ctx context.Context, arg UpdateVerifyAdminEmailParams) error {
	_, err := q.db.ExecContext(ctx, updateVerifyAdminEmail, arg.Email, arg.SecretCode, arg.IsUsed)
	return err
}
