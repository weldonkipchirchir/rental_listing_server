// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"database/sql"
	"time"
)

type Admin struct {
	ID              int32        `json:"id"`
	Username        string       `json:"username"`
	Email           string       `json:"email"`
	PasswordHash    string       `json:"password_hash"`
	CreatedAt       sql.NullTime `json:"created_at"`
	IsEmailVerified bool         `json:"is_email_verified"`
}

type AdminVerifyEmail struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	SecretCode string    `json:"secret_code"`
	IsUsed     bool      `json:"is_used"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}

type Booking struct {
	ID           int32          `json:"id"`
	UserID       int32          `json:"user_id"`
	ListingID    int32          `json:"listing_id"`
	CheckInDate  time.Time      `json:"check_in_date"`
	CheckOutDate time.Time      `json:"check_out_date"`
	TotalAmount  string         `json:"total_amount"`
	Status       sql.NullString `json:"status"`
	CreatedAt    sql.NullTime   `json:"created_at"`
	DeletedAt    sql.NullTime   `json:"deleted_at"`
}

type Favorite struct {
	ID        int32        `json:"id"`
	UserID    int32        `json:"user_id"`
	ListingID int32        `json:"listing_id"`
	CreatedAt sql.NullTime `json:"created_at"`
}

type Listing struct {
	ID          int32          `json:"id"`
	AdminID     int32          `json:"admin_id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	Price       string         `json:"price"`
	Location    sql.NullString `json:"location"`
	CreatedAt   sql.NullTime   `json:"created_at"`
	Available   sql.NullBool   `json:"available"`
	TotalViews  sql.NullInt32  `json:"total_views"`
	Imagelinks  []string       `json:"imagelinks"`
}

type Notification struct {
	ID            int32          `json:"id"`
	UserID        sql.NullInt32  `json:"user_id"`
	Message       string         `json:"message"`
	Read          sql.NullBool   `json:"read"`
	CreatedAt     sql.NullTime   `json:"created_at"`
	Subject       sql.NullString `json:"subject"`
	BookingID     int32          `json:"booking_id"`
	Email         sql.NullString `json:"email"`
	AdminID       sql.NullInt32  `json:"admin_id"`
	SenderAdminID sql.NullInt32  `json:"sender_admin_id"`
	SenderUserID  sql.NullInt32  `json:"sender_user_id"`
}

type Payment struct {
	ID            int32          `json:"id"`
	BookingID     int32          `json:"booking_id"`
	Amount        string         `json:"amount"`
	Status        sql.NullString `json:"status"`
	PaymentMethod sql.NullString `json:"payment_method"`
	TransactionID sql.NullString `json:"transaction_id"`
	PaidAt        sql.NullTime   `json:"paid_at"`
	CreatedAt     sql.NullTime   `json:"created_at"`
	UserID        int32          `json:"user_id"`
}

type Review struct {
	ID        int32          `json:"id"`
	UserID    int32          `json:"user_id"`
	ListingID int32          `json:"listing_id"`
	Rating    int32          `json:"rating"`
	Comment   sql.NullString `json:"comment"`
	CreatedAt sql.NullTime   `json:"created_at"`
}

type Stat struct {
	ID            int32          `json:"id"`
	ListingID     int32          `json:"listing_id"`
	AdminID       int32          `json:"admin_id"`
	TotalViews    sql.NullInt32  `json:"total_views"`
	TotalBookings sql.NullInt32  `json:"total_bookings"`
	AverageRating sql.NullString `json:"average_rating"`
	CreatedAt     sql.NullTime   `json:"created_at"`
}

type User struct {
	ID              int32        `json:"id"`
	Username        string       `json:"username"`
	Email           string       `json:"email"`
	PasswordHash    string       `json:"password_hash"`
	CreatedAt       sql.NullTime `json:"created_at"`
	IsEmailVerified bool         `json:"is_email_verified"`
}

type UserVerifyEmail struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	SecretCode string    `json:"secret_code"`
	IsUsed     bool      `json:"is_used"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}
