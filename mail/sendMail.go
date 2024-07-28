package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error
	SendVerificationEmail(toEmail string, verificationLink string, username string) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewEmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}

func (sender *GmailSender) SendVerificationEmail(toEmail string, verificationLink string, username string) error {
	subject := "Verify Your Email Address"
	content := fmt.Sprintf(`Hello, %s. <br/>
	Thank you for registering with us!<br/>
	Click <a href='%s'>here</a> to verify your email address.`, username, verificationLink)
	return sender.SendEmail(subject, content, []string{toEmail}, nil, nil, nil)
}
