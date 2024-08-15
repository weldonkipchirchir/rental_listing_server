package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/weldonkipchirchir/rental_listing/db/sqlc"
)

func (s *Server) VerifyEmailUser(c *gin.Context) {
	email := c.Param("email")
	token := c.Param("token")

	arg := db.GetVerifyUserEmailParams{
		Email:      email,
		SecretCode: token,
	}

	verify, err := s.q.GetVerifyUserEmail(c, arg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if verify.IsUsed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already verified"})
		return
	}

	err = s.q.UpdateVerifyUserEmail(c, db.UpdateVerifyUserEmailParams{
		Email:      email,
		SecretCode: token,
		IsUsed:     true,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.q.UpdateUserEmailVerified(c, email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Redirect to the login page after successful email verification
	c.Redirect(http.StatusFound, "http://localhost:3000/sign-in")
}

func (s *Server) VerifyEmailAdmin(c *gin.Context) {
	email := c.Param("email")
	token := c.Param("token")

	arg := db.GetVerifyAdminEmailParams{
		Email:      email,
		SecretCode: token,
	}

	verify, err := s.q.GetVerifyAdminEmail(c, arg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if verify.IsUsed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already verified"})
		return
	}

	err = s.q.UpdateVerifyAdminEmail(c, db.UpdateVerifyAdminEmailParams{
		Email:      email,
		SecretCode: token,
		IsUsed:     true,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.q.UpdateAdminEmailVerified(c, email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Redirect to the login page after successful email verification
	c.Redirect(http.StatusFound, "http://localhost:3000/signin-admin")
}
