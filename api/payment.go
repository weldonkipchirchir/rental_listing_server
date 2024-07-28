package api

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	db "github.com/weldonkipchirchir/rental_listing/db/sqlc"
)

type paymentReq struct {
	TotalAmount string `json:"total_amount"`
}

func (s *Server) HandleCreatePaymentIntent(c *gin.Context) {

	var req paymentReq

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	priceFloat, err := strconv.ParseFloat(req.TotalAmount, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Convert to cents and round off
	payAmount := int64(math.Round(priceFloat * 100))

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(payAmount),
		Currency: stripe.String("usd"),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"clientSecret": pi.ClientSecret,
	})
}

// GetPaymentByID retrieves a payment by its ID.
func (s *Server) GetPaymentByID(c *gin.Context) {
	paymentIDStr := c.Param("id")
	paymentID, err := strconv.Atoi(paymentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	payment, err := s.q.GetPaymentByID(c, int32(paymentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// GetPaymentsByBookingID retrieves all payments for a specific booking.
func (s *Server) GetPaymentsByBookingID(c *gin.Context) {
	bookingIDStr := c.Param("booking_id")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	payments, err := s.q.GetPaymentsByBookingID(c, int32(bookingID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// Payment represents a payment entity.
type PaymentUserResponse struct {
	ID            int32     `json:"id"`
	BookingID     int32     `json:"booking_id"`
	Title         string    `json:"title"`
	Amount        string    `json:"amount"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"payment_method"`
	TransactionID string    `json:"transaction_id"`
	PaidAt        time.Time `json:"paid_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// GetPaymentsByUserID retrieves all payments for a specific user.
func (s *Server) GetPaymentsByUserID(c *gin.Context) {
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

	payments, err := s.q.GetPaymentsByUserID(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := make([]PaymentUserResponse, len(payments))
	for i, payment := range payments {
		res[i] = PaymentUserResponse{
			ID:            payment.ID,
			BookingID:     payment.BookingID,
			Title:         payment.Title,
			Amount:        payment.Amount,
			Status:        payment.Status.String,
			PaymentMethod: payment.PaymentMethod.String,
			TransactionID: payment.TransactionID.String,
			PaidAt:        payment.PaidAt.Time,
		}
	}

	c.JSON(http.StatusOK, res)
}

// Payment represents a payment entity.
type PaymentAdminResponse struct {
	ID            int32     `json:"id"`
	BookingID     int32     `json:"booking_id"`
	Title         string    `json:"title"`
	Username      string    `json:"username"`
	Amount        string    `json:"amount"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"payment_method"`
	TransactionID string    `json:"transaction_id"`
	PaidAt        time.Time `json:"paid_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// GetPaymentsByAdminID retrieves all payments for a specific admin.
func (s *Server) GetPaymentsByAdminID(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not found"})
		return
	}

	admin, err := s.q.GetAdmin(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorized admins only"})
		return
	}

	payments, err := s.q.GetPaymentsByAdminID(c, admin.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := make([]PaymentAdminResponse, len(payments))
	for i, payment := range payments {
		res[i] = PaymentAdminResponse{
			ID:            payment.ID,
			BookingID:     payment.BookingID,
			Title:         payment.Title,
			Username:      payment.Username,
			Amount:        payment.Amount,
			Status:        payment.Status.String,
			PaymentMethod: payment.PaymentMethod.String,
			TransactionID: payment.TransactionID.String,
			PaidAt:        payment.PaidAt.Time,
		}
	}

	c.JSON(http.StatusOK, res)
}

type paymentStatusStruct struct {
	Status string `json:"status" binding:"required"`
}

// GetPaymentsByStatus retrieves all payments by a specific status.
func (s *Server) GetPaymentsByStatus(c *gin.Context) {
	var req paymentStatusStruct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payments, err := s.q.GetPaymentsByStatus(c, sql.NullString{String: req.Status, Valid: req.Status != ""})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// UpdatePaymentStatus updates the status of a payment by its ID.
func (s *Server) UpdatePaymentStatus(c *gin.Context) {
	paymentIDStr := c.Param("id")
	paymentID, err := strconv.Atoi(paymentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	var req paymentStatusStruct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdatePaymentStatusParams{
		ID:     int32(paymentID),
		Status: sql.NullString{String: req.Status, Valid: true},
	}

	err = s.q.UpdatePaymentStatus(c, arg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update payment status: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment status updated successfully"})
}

// DeletePayment deletes a payment by its ID.
func (s *Server) DeletePayment(c *gin.Context) {
	paymentIDStr := c.Param("id")
	paymentID, err := strconv.Atoi(paymentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	err = s.q.DeletePayment(c, int32(paymentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment deleted successfully"})
}

func (s *Server) Config(c *gin.Context) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	c.JSON(http.StatusOK, gin.H{
		"publishableKey": os.Getenv("STRIPE_PUBLISHABLE_KEY"),
	})
}
