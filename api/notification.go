package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/weldonkipchirchir/rental_listing/db/sqlc"
)

type createNotificationRequest struct {
	UserID    int32  `json:"user_id" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Subject   string `json:"subject" binding:"required"`
	Email     string `json:"email" binding:"required"`
	BookingID int32  `json:"booking_id" binding:"required"`
}

func (s *Server) CreateNotification(c *gin.Context) {
	var req createNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email not found"})
		return
	}

	admin, err := s.q.GetAdmin(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = s.q.CreateNotification(c, db.CreateNotificationParams{
		UserID:        sql.NullInt32{Int32: req.UserID, Valid: true},
		Message:       req.Message,
		SenderAdminID: sql.NullInt32{Int32: admin.ID, Valid: true},
		Subject:       sql.NullString{String: req.Subject, Valid: true},
		Email:         sql.NullString{String: req.Email, Valid: true},
		BookingID:     req.BookingID,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Notification created successfully"})
}

type createAdminNotificationRequest struct {
	AdminID   int32  `json:"admin_id" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Subject   string `json:"subject" binding:"required"`
	Email     string `json:"email" binding:"required"`
	BookingID int32  `json:"booking_id" binding:"required"`
}

func (s *Server) CreateAdminNotification(c *gin.Context) {
	var req createAdminNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email not found"})
		return
	}

	user, err := s.q.GetUser(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = s.q.CreateAdminNotification(c, db.CreateAdminNotificationParams{
		AdminID:      sql.NullInt32{Int32: req.AdminID, Valid: true},
		Subject:      sql.NullString{String: req.Subject, Valid: true},
		SenderUserID: sql.NullInt32{Int32: user.ID, Valid: true},
		Email:        sql.NullString{String: req.Email, Valid: true},
		BookingID:    req.BookingID,
		Message:      req.Message,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Notification created successfully"})
}

type getAllNotificationsResponse struct {
	ID        int32     `json:"id" binding:"required"`
	UserID    int32     `json:"user_id" binding:"required"`
	AdminID   int32     `json:"admin_id" binding:"required"`
	Message   string    `json:"message" binding:"required"`
	Subject   string    `json:"subject" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	BookingID int32     `json:"booking_id" binding:"required"`
	Read      bool      `json:"read" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Server) GetNotifications(c *gin.Context) {
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

	notifications, err := s.q.GetNotificationsByUserID(c, sql.NullInt32{Int32: user.ID, Valid: true})
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	var response []getAllNotificationsResponse
	for _, notification := range notifications {
		response = append(response, getAllNotificationsResponse{
			ID:        notification.ID,
			UserID:    notification.UserID.Int32,
			AdminID:   notification.ID_2,
			Message:   notification.Message,
			Read:      notification.Read.Bool,
			Subject:   notification.Subject.String,
			Email:     notification.Email,
			BookingID: notification.BookingID,
			CreatedAt: notification.CreatedAt.Time,
		})
	}

	c.JSON(http.StatusOK, response)
}

type getAllAdminNotificationsResponse struct {
	ID        int32     `json:"id" binding:"required"`
	AdminID   int32     `json:"admin_id" binding:"required"`
	UserID    int32     `json:"user_id" binding:"required"`
	Message   string    `json:"message" binding:"required"`
	Subject   string    `json:"subject" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	BookingID int32     `json:"booking_id" binding:"required"`
	Read      bool      `json:"read" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Server) GetAdminNotifications(c *gin.Context) {
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

	notifications, err := s.q.GetNotificationsByAdminID(c, sql.NullInt32{Int32: admin.ID, Valid: true})
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	var response []getAllAdminNotificationsResponse
	for _, notification := range notifications {
		response = append(response, getAllAdminNotificationsResponse{
			ID:        notification.ID,
			AdminID:   notification.AdminID.Int32,
			UserID:    notification.ID_2,
			Message:   notification.Message,
			Read:      notification.Read.Bool,
			Subject:   notification.Subject.String,
			Email:     notification.Email,
			BookingID: notification.BookingID,
			CreatedAt: notification.CreatedAt.Time,
		})
	}

	c.JSON(http.StatusOK, response)
}

type getNotificationsResponse struct {
	UserID    int32     `json:"user_id" binding:"required"`
	Message   string    `json:"message" binding:"required"`
	Subject   string    `json:"subject" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	BookingID int32     `json:"booking_id" binding:"required"`
	Read      bool      `json:"read" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type updateNotificationRequest struct {
	Read *bool `json:"read" binding:"required"` // Using a pointer to handle the presence of the field
}

func (s *Server) UpdateNotification(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var req updateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

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

	arg := db.UpdateNotificationReadStatusParams{
		UserID: sql.NullInt32{Int32: user.ID, Valid: true},
		Read:   sql.NullBool{Bool: *req.Read, Valid: true}, // Dereferencing the pointer to get the boolean value
		ID:     int32(id),
	}

	err = s.q.UpdateNotificationReadStatus(c, arg)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification updated successfully"})
}
func (s *Server) UpdateAdminNotification(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var req updateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

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

	arg := db.UpdateAdminNotificationReadStatusParams{
		AdminID: sql.NullInt32{Int32: admin.ID, Valid: true},
		Read:    sql.NullBool{Bool: *req.Read, Valid: true}, // Dereferencing the pointer to get the boolean value
		ID:      int32(id),
	}

	err = s.q.UpdateAdminNotificationReadStatus(c, arg)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification updated successfully"})
}

func (s *Server) DeleteNotification(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = s.q.DeleteNotification(c, int32(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}

func (s *Server) GetNotificationByID(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

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

	arg := db.GetNotificationByIDParams{
		UserID: sql.NullInt32{Int32: user.ID, Valid: true},
		ID:     int32(id),
	}

	notifications, err := s.q.GetNotificationByID(c, arg)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	response := getNotificationsResponse{
		Message: notifications.Message,
		Read:    notifications.Read.Bool,
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) GetUserUnreadNotifications(c *gin.Context) {
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

	notifications, err := s.q.GetUnreadNotificationsByUserID(c, sql.NullInt32{Int32: user.ID, Valid: true})
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	var response []getNotificationsResponse
	for _, notification := range notifications {
		response = append(response, getNotificationsResponse{
			Message: notification.Message,
			Read:    notification.Read.Bool,
		})
	}

	c.JSON(http.StatusOK, response)
}
func (s *Server) GetAdminUnreadNotifications(c *gin.Context) {
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

	notifications, err := s.q.GetUnreadNotificationsByAdminID(c, sql.NullInt32{Int32: admin.ID, Valid: true})
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	var response []getAllAdminNotificationsResponse
	for _, notification := range notifications {
		response = append(response, getAllAdminNotificationsResponse{
			Message: notification.Message,
			Read:    notification.Read.Bool,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) GetSentAdminNotifications(c *gin.Context) {
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

	notifications, err := s.q.GetSentNotificationsByAdminID(c, sql.NullInt32{Int32: admin.ID, Valid: true})
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	var response []getAllAdminNotificationsResponse
	for _, notification := range notifications {
		response = append(response, getAllAdminNotificationsResponse{
			ID:        notification.ID,
			AdminID:   notification.AdminID.Int32,
			Message:   notification.Message,
			Read:      notification.Read.Bool,
			Subject:   notification.Subject.String,
			Email:     notification.Email,
			BookingID: notification.BookingID,
			CreatedAt: notification.CreatedAt.Time,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) GetSentUserNotifications(c *gin.Context) {
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

	notifications, err := s.q.GetSentNotificationsByUserID(c, sql.NullInt32{Int32: user.ID, Valid: true})
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	var response []getAllNotificationsResponse
	for _, notification := range notifications {
		response = append(response, getAllNotificationsResponse{
			ID:        notification.ID,
			UserID:    notification.UserID.Int32,
			Message:   notification.Message,
			Read:      notification.Read.Bool,
			Subject:   notification.Subject.String,
			Email:     notification.Email,
			BookingID: notification.BookingID,
			CreatedAt: notification.CreatedAt.Time,
		})
	}

	c.JSON(http.StatusOK, response)
}
