package emails

import (
	"bitcoin-service/config"
	"gopkg.in/mail.v2"
)

type Service struct {
	senderEmail string
	dialer      *mail.Dialer
}

func (s Service) SendEmail(receiverEmail, subject, text string) error {
	message := mail.NewMessage()
	message.SetHeader("From", s.senderEmail)
	message.SetHeader("To", receiverEmail)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", text)

	if err := s.dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}

func NewEmailService(cfg *config.AppConfig) *Service {
	dialer := mail.NewDialer(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpUsername, cfg.SmtpPassword)

	return &Service{
		senderEmail: cfg.SmtpUsername,
		dialer:      dialer,
	}
}
