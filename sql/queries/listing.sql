-- name: CreateListing :one
INSERT INTO listings (admin_id, title, description, price, location, available, imageLinks)
VALUES ($1, $2, $3, $4, $5, $6, $7::text[]
)
RETURNING id, admin_id, title, description, price, location, available, imageLinks, created_at;

-- name: GetListingByID :one
SELECT id, admin_id, title, description, price, location, available, imageLinks, created_at
FROM listings
WHERE id = $1;

-- name: GetListings :many
SELECT id, admin_id, title, description, price, location, available, imageLinks, created_at
FROM listings
WHERE available = TRUE
ORDER BY created_at DESC;

-- name: GetAdminListings :many
SELECT id, admin_id, title, description, price, location, available, imageLinks, created_at
FROM listings
WHERE admin_id = $1
ORDER BY created_at DESC;

-- name: GetListingsByAdminID :one
SELECT id, admin_id, title, description, price, location, available, imageLinks, created_at
FROM listings
WHERE admin_id = $1 AND id = $2;

-- name: UpdateListing :exec
UPDATE listings
SET
    title = COALESCE(sqlc.narg(title), title),
    description = COALESCE(sqlc.narg(description), description),
    price = COALESCE(sqlc.narg(price), price),
    location = COALESCE(sqlc.narg(location), location),
    available = COALESCE(sqlc.narg(available), available),
    imageLinks = COALESCE(sqlc.narg(imageLinks), imageLinks)
WHERE id = @id AND admin_id = @admin_id
RETURNING *;

-- name: UpdateTotalViews :exec
UPDATE listings
SET total_views = total_views + 1
WHERE id = $1;

-- name: UpdateListingStatus :exec
UPDATE listings
SET available = $2
WHERE id = $1;

-- name: DeleteListing :exec
DELETE FROM listings
WHERE id = $1 AND admin_id = $2;

-- name: ListingActiveBookingCount :one
SELECT COUNT(*) AS confirmed_count
FROM bookings b
JOIN listings l ON b.listing_id = l.id
WHERE b.listing_id = $1 AND b.status = 'confirmed' AND l.admin_id = $2;

-- name: SearchListings :many
SELECT id, admin_id, title, description, price, location, available, imageLinks, created_at
FROM listings
WHERE 
    (title ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%') 
    AND available = TRUE
ORDER BY created_at DESC;

-- name: Statistics :many
SELECT
    l.id,
    l.title,
    l.total_views,
    COUNT(DISTINCT b.id) AS total_bookings,
    CASE 
        WHEN COUNT(r.id) > 0 THEN ROUND(AVG(r.rating)::numeric, 2)
        ELSE NULL
    END AS average_rating,
    COALESCE(SUM(CASE 
        WHEN b.status IN ('confirmed', 'completed') THEN b.total_amount 
        ELSE 0 
    END), 0) AS total_confirmed_amount
FROM
    listings l
LEFT JOIN
    bookings b ON l.id = b.listing_id
LEFT JOIN
    reviews r ON l.id = r.listing_id
WHERE
    l.admin_id = $1
GROUP BY
    l.id, l.title, l.total_views
ORDER BY
    l.total_views DESC;