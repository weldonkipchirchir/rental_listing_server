package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	db "github.com/weldonkipchirchir/rental_listing/db/sqlc"
	"github.com/weldonkipchirchir/rental_listing/redisCache"
)

// CreateBooking handles the creation of a new booking.
func (s *Server) CreateBooking(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not found"})
		return
	}

	user, err := s.q.GetUser(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorized users only"})
		return
	}

	var req struct {
		ListingID     int       `json:"listing_id" binding:"required"`
		CheckInDate   time.Time `json:"check_in_date" binding:"required"`
		CheckOutDate  time.Time `json:"check_out_date" binding:"required"`
		TotalAmount   string    `json:"total_amount" binding:"required"`
		Status        string    `json:"status" binding:"required"`
		PaymentId     string    `json:"paymentId" binding:"required"`
		PaymentMethod []string  `json:"paymentMethod" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	booking, err := s.q.CreateBooking(c, db.CreateBookingParams{
		UserID:       user.ID,
		ListingID:    int32(req.ListingID),
		CheckInDate:  req.CheckInDate,
		CheckOutDate: req.CheckOutDate,
		TotalAmount:  req.TotalAmount,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	paymentMethod := req.PaymentMethod[0]

	paid_at := time.Now()

	_, err = s.q.CreatePayment(c, db.CreatePaymentParams{
		BookingID:     booking.ID,
		Amount:        booking.TotalAmount,
		Status:        sql.NullString{String: req.Status, Valid: true},
		PaymentMethod: sql.NullString{String: paymentMethod, Valid: true},
		TransactionID: sql.NullString{String: req.PaymentId, Valid: true},
		PaidAt:        sql.NullTime{Time: paid_at, Valid: true},
		UserID:        int32(user.ID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateListingStatusParams{
		ID:        int32(req.ListingID),
		Available: sql.NullBool{Bool: false, Valid: true},
	}
	err = s.q.UpdateListingStatus(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, booking)
}

type GetUserBookingsResponse struct {
	ID            int32  `json:"id"`
	AdminID       int32  `json:"admin_id"`
	AdminUsername string `json:"admin_username"`
	AdminEmail    string `json:"admin_email"`
	Title         string `json:"title"`
	UserID        int32  `json:"user_id"`
	ListingID     int32  `json:"listing_id"`
	CheckInDate   string `json:"check_in_date"`
	CheckOutDate  string `json:"check_out_date"`
	TotalAmount   string `json:"total_amount"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

// GetBookingsByUserID retrieves bookings for the user by email.
func (s *Server) GetBookingsByUserID(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not found"})
		return
	}

	var bookings []GetUserBookingsResponse

	// Fetch from database
	user, err := s.q.GetUser(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized users only"})
		return
	}

	rows, err := s.q.GetBookingsByUserID(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bookings = make([]GetUserBookingsResponse, len(rows))
	for i, row := range rows {
		bookings[i] = GetUserBookingsResponse{
			ID:            row.ID,
			AdminID:       row.ID_2,
			AdminEmail:    row.Email,
			AdminUsername: row.Username,
			Title:         row.Title,
			UserID:        row.UserID,
			ListingID:     row.ListingID,
			CheckInDate:   row.CheckInDate.Format("2006-01-02"),
			CheckOutDate:  row.CheckOutDate.Format("2006-01-02"),
			TotalAmount:   row.TotalAmount,
			Status:        row.Status.String,
			CreatedAt:     row.CreatedAt.Time.Format("2006-01-02"),
		}
	}

	if len(bookings) == 0 {
		bookings = []GetUserBookingsResponse{}
		c.JSON(http.StatusOK, bookings)
	}

	c.JSON(http.StatusOK, bookings)
}

type GetBookingsByAdminResponse struct {
	ID           int32  `json:"id"`
	Title        string `json:"title"`
	Location     string `json:"location"`
	UserID       int32  `json:"user_id"`
	UserUsername string `json:"user_username"`
	UserEmail    string `json:"user_email"`
	ListingID    int32  `json:"listing_id"`
	CheckInDate  string `json:"check_in_date"`
	CheckOutDate string `json:"check_out_date"`
	TotalAmount  string `json:"total_amount"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
}

func (s *Server) GetBookingsByAdminID(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not found"})
		return
	}

	admin, err := s.q.GetAdmin(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorized admin only"})
		return
	}

	rows, err := s.q.GetBookingsByAdminID(c, admin.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Transform database rows into GetBookingsByAdminResponse
	bookings := make([]GetBookingsByAdminResponse, len(rows))
	for i, row := range rows {
		bookings[i] = GetBookingsByAdminResponse{
			ID:           row.ID,
			Title:        row.Title,
			Location:     row.Location.String,
			UserID:       row.UserID,
			UserUsername: row.UserUsername,
			UserEmail:    row.UserEmail,
			ListingID:    row.ListingID,
			CheckInDate:  row.CheckInDate.UTC().Format("2006-01-02"),
			CheckOutDate: row.CheckOutDate.UTC().Format("2006-01-02"),
			TotalAmount:  row.TotalAmount,
			Status:       row.Status.String,
			CreatedAt:    row.CreatedAt.Time.UTC().Format("2006-01-02"),
		}
	}

	c.JSON(http.StatusOK, bookings)
}

// GetBookingByAdminIDAndID retrieves a specific booking for a listing owned by a specific admin.
func (s *Server) GetBookingByAdminIDAndID(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not found"})
		return
	}

	admin, err := s.q.GetAdmin(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized admin only"})
		return
	}

	bookingIDStr := c.Param("id")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	arg := db.GetBookingsByAdminIDByIDParams{
		ID:      int32(bookingID),
		AdminID: admin.ID,
	}

	booking, err := s.q.GetBookingsByAdminIDByID(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, booking)
}

// GetBookingsByListingID retrieves all bookings for a specific listing.
func (s *Server) GetBookingsByListingID(c *gin.Context) {
	listingIDStr := c.Param("id")
	listingID, err := strconv.Atoi(listingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing ID"})
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not found"})
		return
	}

	admin, err := s.q.GetAdmin(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized admin only"})
		return
	}

	cacheKey := fmt.Sprintf("bookingListing:%d", listingID)

	// Try to get data from cache
	var bookings []db.Booking
	err = redisCache.GetCacheRedis(context.Background(), s.redis, cacheKey, &bookings)
	if err == nil && len(bookings) != 0 {
		log.Printf("Retrieved bookings for listing %d from cache", listingID)
		c.JSON(http.StatusOK, bookings)
		return
	} else if err != nil && err != redis.Nil {
		log.Printf("Error fetching bookings for listing %d from cache: %v", listingID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	arg := db.GetBookingsByListingIDParams{
		ListingID: int32(listingID),
		AdminID:   admin.ID,
	}

	rows, err := s.q.GetBookingsByListingID(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, row := range rows {
		booking := db.Booking{
			ID:           row.ID,
			UserID:       row.UserID,
			ListingID:    row.ListingID,
			CheckInDate:  row.CheckInDate,
			CheckOutDate: row.CheckOutDate,
			TotalAmount:  row.TotalAmount,
			Status:       row.Status,
			CreatedAt:    row.CreatedAt,
		}
		bookings = append(bookings, booking)
	}

	// Set cache with expiration
	err = redisCache.SetCacheRedis(context.Background(), s.redis, cacheKey, bookings, time.Hour)
	if err != nil {
		log.Printf("Error setting cache for listing %d: %v", listingID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	log.Printf("Retrieved bookings for listing %d from DB", listingID)

	c.JSON(http.StatusOK, bookings)
}

// UpdateBookingStatus updates the status of a booking by its ID.

// Struct to represent a booking
type Booking struct {
	ID           int
	UserID       int
	ListingID    int
	CheckInDate  time.Time
	CheckOutDate time.Time
	TotalAmount  float64
	Status       string
	CreatedAt    time.Time
}

// Function to handle the update booking status route
func (s *Server) UpdateBookingStatus(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	adminEmail, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found in context"})
		return
	}

	// Fetch admin ID from DB based on email
	admin, err := s.q.GetAdmin(c, adminEmail.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve admin"})
		return
	}

	type updateRequest struct {
		Status string `json:"status,omitempty" binding:"omitempty,min=3"`
	}

	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateBookingStatusByIDAndAdminIDParams{
		ID:      int32(bookingID),
		AdminID: admin.ID,
		Status:  sql.NullString{String: req.Status, Valid: true},
	}

	// Update booking status in DB
	err = s.q.UpdateBookingStatusByIDAndAdminID(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update booking status: %v", err)})
		return
	}

	arg1 := db.GetBookingsByAdminIDByIDParams{
		ID:      int32(bookingID),
		AdminID: admin.ID,
	}
	// Fetch updated booking details from DB
	_, err = s.q.GetBookingsByAdminIDByID(c, arg1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch updated booking: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking status updated successfully"})
}

// DeleteBooking deletes a booking by its ID.
func (s *Server) DeleteBooking(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not found"})
		return
	}

	user, err := s.q.GetUser(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized users only"})
		return
	}

	arg := db.DeleteUserBookingParams{
		ID:     int32(bookingID),
		UserID: user.ID,
	}

	err = s.q.DeleteUserBooking(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	list, err := s.q.GetListingByID(c, int32(bookingID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg1 := db.UpdateListingStatusParams{
		ID:        int32(list.ID),
		Available: sql.NullBool{Bool: true, Valid: true},
	}
	err = s.q.UpdateListingStatus(c, arg1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted successfully"})
}

// upda.
func (s *Server) updateCancelledBooking(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not found"})
		return
	}

	user, err := s.q.GetUser(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized users only"})
		return
	}

	arg := db.UpdateBookingStatusByIDAndUserIDParams{
		ID:     int32(bookingID),
		UserID: user.ID,
	}

	err = s.q.UpdateBookingStatusByIDAndUserID(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	list, err := s.q.GetListingByID(c, int32(bookingID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg1 := db.UpdateListingStatusParams{
		ID:        int32(list.ID),
		Available: sql.NullBool{Bool: true, Valid: true},
	}
	err = s.q.UpdateListingStatus(c, arg1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking updated successfully"})
}
