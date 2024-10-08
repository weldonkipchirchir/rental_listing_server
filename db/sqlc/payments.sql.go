// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: payments.sql

package db

import (
	"context"
	"database/sql"
)

const createPayment = `-- name: CreatePayment :one
INSERT INTO payments (booking_id, amount, status, payment_method, transaction_id, paid_at, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, booking_id, amount, status, payment_method, transaction_id, paid_at, created_at
`

type CreatePaymentParams struct {
	BookingID     int32          `json:"booking_id"`
	Amount        string         `json:"amount"`
	Status        sql.NullString `json:"status"`
	PaymentMethod sql.NullString `json:"payment_method"`
	TransactionID sql.NullString `json:"transaction_id"`
	PaidAt        sql.NullTime   `json:"paid_at"`
	UserID        int32          `json:"user_id"`
}

type CreatePaymentRow struct {
	ID            int32          `json:"id"`
	BookingID     int32          `json:"booking_id"`
	Amount        string         `json:"amount"`
	Status        sql.NullString `json:"status"`
	PaymentMethod sql.NullString `json:"payment_method"`
	TransactionID sql.NullString `json:"transaction_id"`
	PaidAt        sql.NullTime   `json:"paid_at"`
	CreatedAt     sql.NullTime   `json:"created_at"`
}

func (q *Queries) CreatePayment(ctx context.Context, arg CreatePaymentParams) (CreatePaymentRow, error) {
	row := q.db.QueryRowContext(ctx, createPayment,
		arg.BookingID,
		arg.Amount,
		arg.Status,
		arg.PaymentMethod,
		arg.TransactionID,
		arg.PaidAt,
		arg.UserID,
	)
	var i CreatePaymentRow
	err := row.Scan(
		&i.ID,
		&i.BookingID,
		&i.Amount,
		&i.Status,
		&i.PaymentMethod,
		&i.TransactionID,
		&i.PaidAt,
		&i.CreatedAt,
	)
	return i, err
}

const deletePayment = `-- name: DeletePayment :exec
DELETE FROM payments
WHERE id = $1
`

func (q *Queries) DeletePayment(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deletePayment, id)
	return err
}

const getPaymentByID = `-- name: GetPaymentByID :one
SELECT id, booking_id, amount, status, payment_method, transaction_id, paid_at, created_at
FROM payments
WHERE id = $1
`

type GetPaymentByIDRow struct {
	ID            int32          `json:"id"`
	BookingID     int32          `json:"booking_id"`
	Amount        string         `json:"amount"`
	Status        sql.NullString `json:"status"`
	PaymentMethod sql.NullString `json:"payment_method"`
	TransactionID sql.NullString `json:"transaction_id"`
	PaidAt        sql.NullTime   `json:"paid_at"`
	CreatedAt     sql.NullTime   `json:"created_at"`
}

func (q *Queries) GetPaymentByID(ctx context.Context, id int32) (GetPaymentByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getPaymentByID, id)
	var i GetPaymentByIDRow
	err := row.Scan(
		&i.ID,
		&i.BookingID,
		&i.Amount,
		&i.Status,
		&i.PaymentMethod,
		&i.TransactionID,
		&i.PaidAt,
		&i.CreatedAt,
	)
	return i, err
}

const getPaymentsByAdminID = `-- name: GetPaymentsByAdminID :many
SELECT p.id, p.booking_id, l.title, u.username, p.amount, p.status, p.payment_method, p.transaction_id, p.paid_at, p.created_at
FROM payments p
JOIN bookings b ON p.booking_id = b.id
JOIN listings l ON b.listing_id = l.id
JOIN users u ON b.user_id = u.id
WHERE l.admin_id = $1
ORDER BY p.created_at DESC
`

type GetPaymentsByAdminIDRow struct {
	ID            int32          `json:"id"`
	BookingID     int32          `json:"booking_id"`
	Title         string         `json:"title"`
	Username      string         `json:"username"`
	Amount        string         `json:"amount"`
	Status        sql.NullString `json:"status"`
	PaymentMethod sql.NullString `json:"payment_method"`
	TransactionID sql.NullString `json:"transaction_id"`
	PaidAt        sql.NullTime   `json:"paid_at"`
	CreatedAt     sql.NullTime   `json:"created_at"`
}

func (q *Queries) GetPaymentsByAdminID(ctx context.Context, adminID int32) ([]GetPaymentsByAdminIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getPaymentsByAdminID, adminID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPaymentsByAdminIDRow
	for rows.Next() {
		var i GetPaymentsByAdminIDRow
		if err := rows.Scan(
			&i.ID,
			&i.BookingID,
			&i.Title,
			&i.Username,
			&i.Amount,
			&i.Status,
			&i.PaymentMethod,
			&i.TransactionID,
			&i.PaidAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPaymentsByBookingID = `-- name: GetPaymentsByBookingID :many
SELECT id, booking_id, amount, status, payment_method, transaction_id, paid_at, created_at
FROM payments
WHERE booking_id = $1
ORDER BY created_at DESC
`

type GetPaymentsByBookingIDRow struct {
	ID            int32          `json:"id"`
	BookingID     int32          `json:"booking_id"`
	Amount        string         `json:"amount"`
	Status        sql.NullString `json:"status"`
	PaymentMethod sql.NullString `json:"payment_method"`
	TransactionID sql.NullString `json:"transaction_id"`
	PaidAt        sql.NullTime   `json:"paid_at"`
	CreatedAt     sql.NullTime   `json:"created_at"`
}

func (q *Queries) GetPaymentsByBookingID(ctx context.Context, bookingID int32) ([]GetPaymentsByBookingIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getPaymentsByBookingID, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPaymentsByBookingIDRow
	for rows.Next() {
		var i GetPaymentsByBookingIDRow
		if err := rows.Scan(
			&i.ID,
			&i.BookingID,
			&i.Amount,
			&i.Status,
			&i.PaymentMethod,
			&i.TransactionID,
			&i.PaidAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPaymentsByStatus = `-- name: GetPaymentsByStatus :many
SELECT id, booking_id, amount, status, payment_method, transaction_id, paid_at, created_at
FROM payments
WHERE status = $1
ORDER BY created_at DESC
`

type GetPaymentsByStatusRow struct {
	ID            int32          `json:"id"`
	BookingID     int32          `json:"booking_id"`
	Amount        string         `json:"amount"`
	Status        sql.NullString `json:"status"`
	PaymentMethod sql.NullString `json:"payment_method"`
	TransactionID sql.NullString `json:"transaction_id"`
	PaidAt        sql.NullTime   `json:"paid_at"`
	CreatedAt     sql.NullTime   `json:"created_at"`
}

func (q *Queries) GetPaymentsByStatus(ctx context.Context, status sql.NullString) ([]GetPaymentsByStatusRow, error) {
	rows, err := q.db.QueryContext(ctx, getPaymentsByStatus, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPaymentsByStatusRow
	for rows.Next() {
		var i GetPaymentsByStatusRow
		if err := rows.Scan(
			&i.ID,
			&i.BookingID,
			&i.Amount,
			&i.Status,
			&i.PaymentMethod,
			&i.TransactionID,
			&i.PaidAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPaymentsByUserID = `-- name: GetPaymentsByUserID :many
SELECT p.id, p.booking_id, p.amount, l.title, p.status, p.payment_method, p.transaction_id, p.paid_at, p.created_at
FROM payments p
JOIN bookings b ON p.booking_id = b.id
JOIN listings l ON b.listing_id = l.id
WHERE b.user_id = $1
ORDER BY p.created_at DESC
`

type GetPaymentsByUserIDRow struct {
	ID            int32          `json:"id"`
	BookingID     int32          `json:"booking_id"`
	Amount        string         `json:"amount"`
	Title         string         `json:"title"`
	Status        sql.NullString `json:"status"`
	PaymentMethod sql.NullString `json:"payment_method"`
	TransactionID sql.NullString `json:"transaction_id"`
	PaidAt        sql.NullTime   `json:"paid_at"`
	CreatedAt     sql.NullTime   `json:"created_at"`
}

func (q *Queries) GetPaymentsByUserID(ctx context.Context, userID int32) ([]GetPaymentsByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getPaymentsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPaymentsByUserIDRow
	for rows.Next() {
		var i GetPaymentsByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.BookingID,
			&i.Amount,
			&i.Title,
			&i.Status,
			&i.PaymentMethod,
			&i.TransactionID,
			&i.PaidAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePaymentStatus = `-- name: UpdatePaymentStatus :exec
UPDATE payments
SET status = COALESCE($2, status)
WHERE id = $1
RETURNING id, booking_id, amount, status, payment_method, transaction_id, paid_at, created_at
`

type UpdatePaymentStatusParams struct {
	ID     int32          `json:"id"`
	Status sql.NullString `json:"status"`
}

func (q *Queries) UpdatePaymentStatus(ctx context.Context, arg UpdatePaymentStatusParams) error {
	_, err := q.db.ExecContext(ctx, updatePaymentStatus, arg.ID, arg.Status)
	return err
}
