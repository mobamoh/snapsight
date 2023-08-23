package models

import (
	"fmt"
	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "mobamoh@snapsight.de"
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
	//return &EmailService{
	//	dialer: mail.NewDialer("sandbox.smtp.mailtrap.io",
	//		587,
	//		"50f04ed0de1c7f",
	//		"2d5d51f442391c"),
	//}
}

func (service *EmailService) Send(email Email) error {
	msg := mail.NewMessage()
	msg.SetHeader("To", email.To)
	service.setFrom(msg, email)
	msg.SetHeader("Subject", email.Subject)
	service.setBody(msg, email)
	err := service.dialer.DialAndSend()
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func (service *EmailService) setFrom(msg *mail.Message, email Email) {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case service.DefaultSender != "":
		from = service.DefaultSender
	default:
		from = DefaultSender
	}
	msg.SetHeader("From", from)
}

func (service *EmailService) setBody(msg *mail.Message, email Email) {
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

func (service *EmailService) ForgotPassword(to, resetURL string) error {
	email := Email{
		Subject:   "Reset your password",
		To:        to,
		Plaintext: "To reset your password, please visit the following link: " + resetURL,
		HTML:      `<p>To reset your password, please visit the following link: <a href="` + resetURL + `">` + resetURL + `</a></p>`,
	}

	err := service.Send(email)
	if err != nil {
		return fmt.Errorf("email forgot password: %w", err)
	}
	return nil
}
