package api

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/weldonkipchirchir/rental_listing/db/sqlc"
	"github.com/weldonkipchirchir/rental_listing/util"
)

var keyHex = "2b7e151628aed2a6abf7158809cf4f3c2b7e151628aed2a6abf7158809cf4f3c"

func (s *Server) VerifyEmailUser(c *gin.Context) {
	encryptEMail := c.Param("email")
	encryptToken := c.Param("token")

	// Convert hexadecimal string to byte slice
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		fmt.Println("Error decoding hex string:", err)
		return
	}

	decryptedEmailBytes, err := util.Decrypt(encryptEMail, key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email := string(decryptedEmailBytes)

	decryptedCodeBytes, err := util.Decrypt(encryptToken, key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token := string(decryptedCodeBytes)

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
	c.Redirect(http.StatusFound, "http://localhost:5173/sign-in")
}

func (s *Server) VerifyEmailAdmin(c *gin.Context) {
	email := c.Param("email")
	token := c.Param("token")

	// Convert hexadecimal string to byte slice
	// key, err := hex.DecodeString(keyHex)
	// if err != nil {
	// 	fmt.Println("Error decoding hex string:", err)
	// 	return
	// }

	// decryptedEmailBytes, err := util.Decrypt(encryptEMail, key)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// email := string(decryptedEmailBytes)

	// decryptedCodeBytes, err := util.Decrypt(encryptToken, key)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// token := string(decryptedCodeBytes)

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
	c.Redirect(http.StatusFound, "http://localhost:5173/signin-admin")
}
