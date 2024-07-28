-- name: CreateBooking :one
INSERT INTO bookings (user_id, listing_id, check_in_date, check_out_date, total_amount)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, listing_id, check_in_date, check_out_date, total_amount, status, created_at;

-- name: GetUserBookingByID :one
SELECT id, user_id, listing_id, check_in_date, check_out_date, total_amount, status, created_at
FROM bookings
WHERE id = $1 AND user_id = $2;

-- name: GetBookingsByUserID :many
SELECT b.id, l.title, a.id, a.username, a.email, b.user_id, b.listing_id, b.check_in_date, b.check_out_date, b.total_amount, b.status, b.created_at
FROM bookings b
JOIN listings l ON b.listing_id = l.id
JOIN admins a ON l.admin_id = a.id
WHERE b.user_id = $1 AND deleted_at IS NULL
ORDER BY b.created_at DESC;

-- name: GetBookingsByAdminID :many
SELECT 
    b.id, 
    l.title,
    l.location, 
    b.user_id, 
    u.username AS user_username,
    u.email AS user_email,
    b.listing_id, 
    b.check_in_date, 
    b.check_out_date, 
    b.total_amount, 
    b.status, 
    b.created_at
FROM 
    bookings b
JOIN 
    listings l ON b.listing_id = l.id
JOIN 
    users u ON b.user_id = u.id  
WHERE 
    l.admin_id = $1 AND (b.status = 'completed' OR b.status = 'confirmed' OR b.status = 'pending')
ORDER BY 
    b.created_at DESC;


-- name: GetBookingsByAdminIDByID :one
SELECT b.id, b.user_id, b.listing_id, b.check_in_date, b.check_out_date, b.total_amount, b.status, b.created_at
FROM bookings b
JOIN listings l ON b.listing_id = l.id
WHERE l.admin_id = $1 AND b.id = $2
ORDER BY b.created_at DESC;

-- name: GetBookingsByListingID :many
SELECT b.id, b.user_id, b.listing_id, b.check_in_date, b.check_out_date, b.total_amount, b.status, b.created_at
FROM bookings b
JOIN listings l ON b.listing_id = l.id
WHERE b.listing_id = $1 AND l.admin_id = $2
ORDER BY b.created_at DESC;

-- name: UpdateBookingStatusByIDAndAdminID :exec
UPDATE bookings b
SET status = COALESCE(sqlc.narg(status), status)
FROM listings l
WHERE b.listing_id = l.id AND b.id = @id AND l.admin_id = @admin_id
RETURNING b.id, b.user_id, b.listing_id, b.check_in_date, b.check_out_date, b.total_amount, b.status, b.created_at;

-- name: UpdateBookingStatusByIDAndUserID :exec
UPDATE bookings
SET status = 'cancelled'
WHERE id = $1 AND user_id = $2;

-- name: DeleteUserBooking :exec
UPDATE bookings
SET deleted_at = NOW()
WHERE id = $1 AND user_id = $2;

-- name: CountBookingsByUserID :one
SELECT COUNT(*)
FROM bookings
WHERE user_id = $1;


