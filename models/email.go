package models

import (
	"fmt"
	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "mo@bamoh.de"
)

type Email struct {
	From      string
	To        string
	Subject   string
	Plaintext string
	HTML      string
}

type EmailService struct {
	DefaultSender string
	dialer        *mail.Dialer
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(config SMTPConfig) *EmailService {
	return &EmailService{
		dialer: mail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}
}

func (r *EmailService) Send(email Email) error {
	msg := mail.NewMessage()
	msg.SetHeader("To", email.To)
	r.setFrom(msg, email)
	msg.SetHeader("From", email.From)
	msg.SetHeader("Subject", email.Subject)
	r.setBody(msg, email)

	err := r.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func (r *EmailService) setFrom(msg *mail.Message, email Email) {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case r.DefaultSender != "":
		from = r.DefaultSender
	default:
		from = DefaultSender
	}
	msg.SetHeader("From", from)
}

func (r *EmailService) setBody(msg *mail.Message, email Email) {
	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBody("text/plain", email.Plaintext)
		msg.AddAlternative("text/html", email.HTML)
	case email.Plaintext != "":
		msg.SetBody("text/plain", email.Plaintext)
	case email.HTML != "":
		msg.SetBody("text/html", email.HTML)
	}
}

func (r *EmailService) ForgotPassword(to, resetURL string) error {
	email := Email{
		Subject:   "Reset your password",
		To:        to,
		Plaintext: "To reset your password, please visit the following link: " + resetURL,
		HTML:      `<p>To reset your password, please visit the following link: <a href="` + resetURL + `">` + resetURL + `</a></p>`,
	}

	err := r.Send(email)
	if err != nil {
		return fmt.Errorf("email forgot password: %w", err)
	}
	return nil
}
