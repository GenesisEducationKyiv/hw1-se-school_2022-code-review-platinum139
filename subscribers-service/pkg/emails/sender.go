package emails

import (
	"gopkg.in/mail.v2"
	"subscribers-service/config"
)

type Service struct {
	senderEmail string
	dialer      *mail.Dialer
}

func (s Service) Send(receiverEmail, subject, text string) error {
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
	dialer := mail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword)

	return &Service{
		senderEmail: cfg.SMTPUsername,
		dialer:      dialer,
	}
}
