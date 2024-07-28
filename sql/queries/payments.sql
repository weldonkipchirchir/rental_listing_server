-- name: CreatePayment :one
INSERT INTO payments (booking_id, amount, status, payment_method, transaction_id, paid_at, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, booking_id, amount, status, payment_method, transaction_id, paid_at, created_at;

-- name: GetPaymentByID :one
SELECT id, booking_id, amount, status, payment_method, transaction_id, paid_at, created_at
FROM payments
WHERE id = $1;

-- name: GetPaymentsByBookingID :many
SELECT id, booking_id, amount, status, payment_method, transaction_id, paid_at, created_at
FROM payments
WHERE booking_id = $1
ORDER BY created_at DESC;

-- name: GetPaymentsByUserID :many
SELECT p.id, p.booking_id, p.amount, l.title, p.status, p.payment_method, p.transaction_id, p.paid_at, p.created_at
FROM payments p
JOIN bookings b ON p.booking_id = b.id
JOIN listings l ON b.listing_id = l.id
WHERE b.user_id = $1
ORDER BY p.created_at DESC;

-- name: GetPaymentsByAdminID :many
SELECT p.id, p.booking_id, l.title, u.username, p.amount, p.status, p.payment_method, p.transaction_id, p.paid_at, p.created_at
FROM payments p
JOIN bookings b ON p.booking_id = b.id
JOIN listings l ON b.listing_id = l.id
JOIN users u ON b.user_id = u.id
WHERE l.admin_id = $1
ORDER BY p.created_at DESC;

-- name: GetPaymentsByStatus :many
SELECT id, booking_id, amount, status, payment_method, transaction_id, paid_at, created_at
FROM payments
WHERE status = $1
ORDER BY created_at DESC;

-- name: UpdatePaymentStatus :exec
UPDATE payments
SET status = COALESCE(sqlc.arg(status), status)
WHERE id = $1
RETURNING id, booking_id, amount, status, payment_method, transaction_id, paid_at, created_at;

-- name: DeletePayment :exec
DELETE FROM payments
WHERE id = $1;
