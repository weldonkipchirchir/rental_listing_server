package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	db "github.com/weldonkipchirchir/rental_listing/db/sqlc"
	"github.com/weldonkipchirchir/rental_listing/tasks"
	"github.com/weldonkipchirchir/rental_listing/token"
	"github.com/weldonkipchirchir/rental_listing/util"
)

type createAdminRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}
type adminResponse struct {
	ID        int32     `json:"id" binding:"required,alphanum"`
	Username  string    `json:"username" binding:"required,alphanum"`
	Email     string    `json:"email" binding:"required,email"`
	CreatedAt time.Time `json:"created_at"`
}

func newAdminResponse(admin db.CreateAdminRow) adminResponse {
	return adminResponse{
		ID:        admin.ID,
		Username:  admin.Username,
		Email:     admin.Email,
		CreatedAt: admin.CreatedAt.Time,
	}
}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (s *Server) CreateAdmin(c *gin.Context) {
	var request createAdminRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := s.q.GetUser(c, request.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}
	if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = s.q.GetAdmin(c, request.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}
	if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if hashedPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to hash password"})
		return
	}

	arg := db.CreateAdminParams{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: hashedPassword,
	}

	admin, err := s.q.CreateAdmin(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	code := util.RandomInt(10000, 99999)
	secretCode := strconv.Itoa(int(code))
	verificationLink := fmt.Sprintf("http://localhost:8000/api/admin/verify/%s/%s", request.Email, secretCode)
	log.Printf("Generated verification link: %s", verificationLink)
	task, err := tasks.NewVerificationEmailTask(admin.Email, verificationLink, admin.Username)
	if err != nil {
		log.Printf("Failed to create task: %v", err)
	}
	info, err := s.client.Enqueue(task)
	if err != nil {
		log.Printf("Failed to enqueue task: %v", err)
	}
	log.Printf("Enqueued task: id=%s queue=%s", info.ID, info.Queue)

	arg1 := db.CreateVerifyAdminEmailParams{
		Username:   request.Username,
		Email:      request.Email,
		SecretCode: secretCode,
	}

	_, err = s.q.CreateVerifyAdminEmail(c, arg1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newAdminResponse(admin)
	c.JSON(http.StatusCreated, rsp)
}

type loginAdminRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginAdminResponse struct {
	Email        string `json:"email" binding:"required,email"`
	Username     string `json:"username" binding:"required,alphanum"`
	Role         string `json:"role" binding:"required,alphanum"`
	Token        string `json:"token" binding:"required,alphanum"`
	RefreshToken string `json:"refreshToken" binding:"required,alphanum"`
}

// SetCookie sets a cookie with the SameSite attribute.
func SetCookie(c *gin.Context, name, value string, maxAge int, path, domain string, secure, httpOnly bool, sameSite http.SameSite) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: sameSite,
		Expires:  time.Now().Add(time.Duration(maxAge) * time.Second),
	}

	http.SetCookie(c.Writer, cookie)
}

func (s *Server) loginAdmin(c *gin.Context) {
	var request loginAdminRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	admin, err := s.q.GetAdmin(c, request.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !admin.IsEmailVerified {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = util.VerifyPassword(request.Password, admin.PasswordHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	token, refreshToken, err := token.GenerateJwtToken(admin.ID, admin.Email, admin.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	SetCookie(c, "token", token, 300*3600, "/", "localhost", false, true, http.SameSiteLaxMode)
	SetCookie(c, "refreshToken", refreshToken, 300*3600, "/", "localhost", false, true, http.SameSiteLaxMode)

	adminRes := loginAdminResponse{
		Username:     admin.Username,
		Email:        admin.Email,
		Role:         "admin",
		Token:        token,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, gin.H{"user": adminRes})
}

type updateAdminRequest struct {
	CurrentPassword string `json:"currentPassword,omitempty" binding:"omitempty,min=6"`
	NewPassword     string `json:"newPassword,omitempty" binding:"omitempty,min=6"`
	Username        string `json:"username,omitempty" binding:"omitempty,min=3"`
}

func (s *Server) updateAdmin(c *gin.Context) {
	var req updateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var errMsg string
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				switch e.Field() {
				case "CurrentPassword":
					errMsg = "Current password must be at least 6 characters long"
				case "NewPassword":
					errMsg = "New password must be at least 6 characters long"
				case "Username":
					errMsg = "Username must be at least 3 characters long"
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
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

	if req.NewPassword == "" || req.Username == "" || req.CurrentPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please enter all fields"})
		return
	}

	err = util.VerifyPassword(req.CurrentPassword, admin.PasswordHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid current email or password"})
		return
	}

	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "internal server error"})
		return
	}

	arg := db.UpdateAdminPasswordParams{
		ID:           admin.ID,
		PasswordHash: hashedPassword,
		Username:     req.Username,
	}

	err = s.q.UpdateAdminPassword(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully updated"})
}

func (s *Server) logoutAdmin(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.SetCookie("refreshToken", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

type forgotAdminPasswordReq struct {
	Email string `json:"email" binding:"required,email"`
}

func (s *Server) forgotAdminPassword(c *gin.Context) {
	var req forgotAdminPasswordReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, err := s.q.GetAdmin(c, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "admin not found"})
		return
	}

	password := util.RandomChars(20)

	hashPassword, err := util.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	arg := db.UpdateAdminForgotPasswordParams{
		PasswordHash: hashPassword,
		Email:        admin.Email,
	}

	err = s.q.UpdateAdminForgotPassword(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add task to send new password email
	task, err := tasks.NewForgotPasswordEmailTask(req.Email, password, admin.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = s.client.Enqueue(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "new password sent to email"})
}
