-- name: CreateFavorite :one
INSERT INTO favorites (user_id, listing_id)
VALUES ($1, $2)
RETURNING id, user_id, listing_id, created_at;

-- name: GetFavoriteByID :one
SELECT id, user_id, listing_id, created_at
FROM favorites
WHERE id = $1 and user_id = $2;

-- name: GetFavoriteByListingID :one
SELECT id, user_id, listing_id, created_at
FROM favorites
WHERE listing_id = $1 and user_id = $2;

-- name: GetListingFavoriteByUser :one
SELECT id, user_id, listing_id, created_at
FROM favorites
WHERE listing_id = $1 and user_id = $2;

-- name: GetFavorite :many
SELECT id, user_id, listing_id, created_at
FROM favorites
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteFavorite :exec
DELETE FROM favorites
WHERE listing_id = $1 and user_id = $2;

-- name: SearchFavorite :many
SELECT f.id, l.admin_id, l.title, l.description, l.price, l.location, l.available, l.imageLinks
FROM favorites f
JOIN listings l ON f.listing_id = l.id
WHERE 
    (l.title ILIKE '%' || $1 || '%' OR l.description ILIKE '%' || $1 || '%') 
    AND l.available = TRUE AND f.user_id = $2
ORDER BY f.created_at DESC;