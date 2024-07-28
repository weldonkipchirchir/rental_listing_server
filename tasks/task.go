package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/weldonkipchirchir/rental_listing/mail"
)

const (
	TypeVerificationEmail   = "email:verification"
	TypeForgotPasswordEmail = "email:forgot_password" // New task type
)

type VerificationEmailPayload struct {
	ToEmail          string `json:"to_email"`
	VerificationLink string `json:"verification_link"`
	Username         string `json:"username"`
}

func NewVerificationEmailTask(toEmail string, verificationLink string, username string) (*asynq.Task, error) {
	payload, err := json.Marshal(VerificationEmailPayload{
		ToEmail:          toEmail,
		VerificationLink: verificationLink,
		Username:         username,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeVerificationEmail, payload), nil
}

func HandleVerificationEmailTask(ctx context.Context, t *asynq.Task) error {
	var payload VerificationEmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	// Send verification email using your EmailSender implementation
	sender := mail.NewEmailSender("Rental Listing", "weldonkipchirchir23@gmail.com", "ukeodaruubrozovp")
	if err := sender.SendVerificationEmail(payload.ToEmail, payload.VerificationLink, payload.Username); err != nil {
		return err
	}

	log.Printf("Sent verification email to: %s", payload.ToEmail)
	return nil
}

type ForgotPasswordEmailPayload struct {
	ToEmail     string `json:"to_email"`
	NewPassword string `json:"new_password"`
	Username    string `json:"username"`
}

func NewForgotPasswordEmailTask(toEmail, newPassword, username string) (*asynq.Task, error) {
	payload, err := json.Marshal(ForgotPasswordEmailPayload{
		ToEmail:     toEmail,
		NewPassword: newPassword,
		Username:    username,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeForgotPasswordEmail, payload), nil
}

func HandleForgotPasswordEmailTask(ctx context.Context, t *asynq.Task) error {
	var payload ForgotPasswordEmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	email := os.Getenv("email")
	password := os.Getenv("passwordEmail")
	// Send forgot password email using your EmailSender implementation
	sender := mail.NewEmailSender("Rental Listing", email, password)
	content := fmt.Sprintf("Hello %s, your new password is: %s. Please change the password", payload.Username, payload.NewPassword)
	if err := sender.SendEmail("Your New Password", content, []string{payload.ToEmail}, nil, nil, nil); err != nil {
		return err
	}

	log.Printf("Sent forgot password email to: %s", payload.ToEmail)
	return nil
}
