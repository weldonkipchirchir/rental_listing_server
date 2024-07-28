-- name: CreateReview :one
INSERT INTO reviews (user_id, listing_id, rating, comment)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, listing_id, rating, comment, created_at;

-- name: GetReviewByID :one
SELECT id, user_id, listing_id, rating, comment, created_at
FROM reviews
WHERE id = $1 and user_id = $2;

-- name: GetReviewsByListingID :many
SELECT r.id, r.user_id, r.listing_id, r.rating, r.comment, r.created_at
FROM reviews r
JOIN listings l ON r.listing_id = l.id
WHERE listing_id = $1 AND admin_id= $2;

-- name: GetListingReviews :many
SELECT r.id, r.user_id, r.listing_id, r.rating, r.comment, r.created_at, u.username
FROM reviews r
JOIN listings l ON r.listing_id = l.id
JOIN users u ON r.user_id = u.id
WHERE listing_id = $1;

-- name: DeleteReview :exec
DELETE FROM reviews
WHERE id = $1 and user_id = $2;