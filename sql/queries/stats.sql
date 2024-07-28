-- name: CreateStat :one
INSERT INTO stats (listing_id, total_views, total_bookings, average_rating)
VALUES ($1, $2, $3, $4)
RETURNING id, listing_id, total_views, total_bookings, average_rating, created_at;

-- name: GetStatByID :one
SELECT id, listing_id, total_views, total_bookings, average_rating, created_at
FROM stats
WHERE id = $1;

-- name: GetStatsByListingID :one
SELECT id, listing_id, total_views, total_bookings, average_rating, created_at
FROM stats
WHERE listing_id = $1;

-- name: UpdateStat :exec
UPDATE stats
SET
    total_views = $2,
    total_bookings = $3,
    average_rating = $4
WHERE id = $1
RETURNING id, listing_id, total_views, total_bookings, average_rating, created_at;

-- name: DeleteStat :exec
DELETE FROM stats
WHERE id = $1;

-- name: GetAllStats :many
SELECT id, listing_id, total_views, total_bookings, average_rating, created_at
FROM stats
ORDER BY created_at DESC;
