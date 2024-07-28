-- name: CreateNotification :one
INSERT INTO notifications (user_id, subject, sender_admin_id, email, booking_id, message)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CreateAdminNotification :one
INSERT INTO notifications (admin_id, subject, sender_user_id, email, booking_id, message)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetNotificationByID :one
SELECT id, user_id, message, read, email, created_at
FROM notifications
WHERE id = $1 AND user_id = $2;

-- name: GetNotificationsByUserID :many
SELECT n.id, n.user_id, a.id, n.subject, n.booking_id, a.email, n.message, n.read, n.created_at
FROM notifications n
JOIN admins a ON n.sender_admin_id = a.id
WHERE n.user_id = $1
ORDER BY n.created_at DESC;

-- name: GetSentNotificationsByUserID :many
SELECT n.id, n.user_id, n.subject, n.booking_id, a.email, n.message, n.read, n.created_at
FROM notifications n
JOIN admins a ON n.admin_id = a.id
WHERE n.sender_user_id = $1
ORDER BY n.created_at DESC;

-- name: GetNotificationsByAdminID :many
SELECT n.id, n.admin_id, a.id, n.subject, n.booking_id, a.email, n.message, n.read, n.created_at
FROM notifications n
JOIN users a ON n.sender_user_id = a.id
WHERE n.admin_id = $1
ORDER BY n.created_at DESC;

-- name: GetSentNotificationsByAdminID :many
SELECT n.id, n.admin_id, n.subject, n.booking_id, a.email, n.message, n.read, n.created_at
FROM notifications n
JOIN users a ON n.user_id = a.id
WHERE n.sender_admin_id = $1
ORDER BY n.created_at DESC;

-- name: GetUnreadNotificationsByUserID :many
SELECT id, user_id, message, read, email, created_at
FROM notifications
WHERE user_id = $1 AND read = FALSE
ORDER BY created_at DESC;

-- name: GetUnreadNotificationsByAdminID :many
SELECT id, admin_id, message, read, email, created_at
FROM notifications
WHERE admin_id = $1 AND read = FALSE
ORDER BY created_at DESC;

-- name: UpdateNotificationReadStatus :exec
UPDATE notifications
SET read = COALESCE(sqlc.arg(read), read)
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, message, read, created_at;

-- name: UpdateAdminNotificationReadStatus :exec
UPDATE notifications
SET read = COALESCE(sqlc.arg(read), read)
WHERE id = $1 AND admin_id = $2
RETURNING id, admin_id, message, read, created_at;

-- name: DeleteNotification :exec
DELETE FROM notifications
WHERE id = $1;

