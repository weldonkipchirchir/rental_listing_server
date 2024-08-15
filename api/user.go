package api

import (
	"database/sql"
	"fmt"
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

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}
type userResponse struct {
	ID        int32     `json:"id" binding:"required,alphanum"`
	Username  string    `json:"username" binding:"required,alphanum"`
	Email     string    `json:"email" binding:"required,email"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.CreateUserRow) userResponse {
	return userResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
	}
}

func (s *Server) CreateUser(c *gin.Context) {
	var request createUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check if the user already exists
	_, err := s.q.GetAdmin(c, request.Email)
	if err == nil {
		// User exists
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}
	if err != sql.ErrNoRows {
		// Some other error occurred
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = s.q.GetUser(c, request.Email)
	if err == nil {
		// User exists
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}
	if err != sql.ErrNoRows {
		// Some other error occurred
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

	arg := db.CreateUserParams{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: hashedPassword,
	}

	user, err := s.q.CreateUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	code := util.RandomInt(10000, 99999)
	secretCode := strconv.Itoa(int(code))

	// Generate verification link (this can be more secure, e.g., using JWT)
	verificationLink := fmt.Sprintf("http://localhost:8000/api/user/verify/%s/%s", request.Email, secretCode)

	// Enqueue verification email task
	task, err := tasks.NewVerificationEmailTask(user.Email, verificationLink, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = s.client.Enqueue(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg1 := db.CreateVerifyUserEmailParams{
		Username:   request.Username,
		Email:      request.Email,
		SecretCode: secretCode,
	}

	_, err = s.q.CreateVerifyUserEmail(c, arg1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	c.JSON(http.StatusCreated, rsp)
}

type loginUserequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
type loginUserResponse struct {
	Email        string `json:"email" binding:"required,email"`
	Username     string `json:"username" binding:"required,alphanum"`
	Role         string `json:"role" binding:"required,alphanum"`
	Token        string `json:"token" binding:"required,alphanum"`
	RefreshToken string `json:"refreshToken" binding:"required,alphanum"`
}

func (s *Server) loginUser(c *gin.Context) {
	var request loginUserequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.q.GetUser(c, request.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !user.IsEmailVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not verified"})
		return
	}

	err = util.VerifyPassword(request.Password, user.PasswordHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect email or password"})
		return
	}

	token, refreshToken, err := token.GenerateJwtToken(user.ID, user.Email, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	userRes := loginUserResponse{
		Username:     user.Username,
		Email:        user.Email,
		Role:         "user",
		Token:        token,
		RefreshToken: refreshToken,
	}

	host := c.Request.Host

	SetCookie(c, "token", token, 300*3600, "/", host, false, true, http.SameSiteLaxMode)
	SetCookie(c, "refreshToken", refreshToken, 300*3600, "/", host, false, true, http.SameSiteLaxMode)

	c.JSON(http.StatusOK, gin.H{"user": userRes})
}

type updateUserRequest struct {
	CurrentPassword string `json:"currentPassword,omitempty" binding:"omitempty,min=6"`
	NewPassword     string `json:"newPassword,omitempty" binding:"omitempty,min=6"`
	Username        string `json:"username,omitempty" binding:"omitempty,min=3"`
}

func (s *Server) updateUser(c *gin.Context) {
	var req updateUserRequest
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

	user, err := s.q.GetUser(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorized user only"})
		return
	}

	if req.NewPassword == "" || req.Username == "" || req.CurrentPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please enter all fields"})
		return
	}

	err = util.VerifyPassword(req.CurrentPassword, user.PasswordHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid current password"})
		return
	}

	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "internal server error"})
		return
	}

	arg := db.UpdateUserPasswordParams{
		ID:           user.ID,
		PasswordHash: hashedPassword,
		Username:     req.Username,
	}

	err = s.q.UpdateUserPassword(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully updated"})
}
func (s *Server) logoutUser(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.SetCookie("refreshToken", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

type forgotPasswordReq struct {
	Email string `json:"email" binding:"required,email"`
}

func (s *Server) forgotPassword(c *gin.Context) {
	var req forgotPasswordReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.q.GetUser(c, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	password := util.RandomChars(20)

	hashPassword, err := util.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	arg := db.UpdateUserForgotPasswordParams{
		PasswordHash: hashPassword,
		Email:        user.Email,
	}

	err = s.q.UpdateUserForgotPassword(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add task to send new password email
	task, err := tasks.NewForgotPasswordEmailTask(req.Email, password, user.Username)
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
