package subscribers

import (
	"fmt"
	"log"
)

type Storage interface {
	Add(record string) error
	GetAll() ([]string, error)
}

type EmailSender interface {
	SendEmail(receiverEmail string, subject string, text string) error
}

type Service struct {
	log        *log.Logger
	storage    Storage
	mailSender EmailSender
}

func (s Service) Add(subscriber string) error {
	return s.storage.Add(subscriber)
}

func (s Service) SendEmails(rate float64, fromCurrency, toCurrency string) error {
	subscribers, err := s.getAll()
	if err != nil {
		s.log.Print("Unable to get subscribers from storage:", err)
		return err
	}

	var failedEmails []string
	for _, subscriber := range subscribers {
		message := fmt.Sprintf("Rate = %.2f", rate)
		subject := fmt.Sprintf("%s to %s rate", fromCurrency, toCurrency)

		err := s.mailSender.SendEmail(subscriber, subject, message)
		if err != nil {
			s.log.Printf("Unable to send mails via mail service for %s: %s", subscriber, err)
			failedEmails = append(failedEmails, subscriber)
		}
	}

	if len(failedEmails) != 0 {
		return SendMailError{Subscribers: failedEmails}
	}

	return nil
}

func (s Service) getAll() ([]string, error) {
	return s.storage.GetAll()
}

func NewSubscribersService(logger *log.Logger, storage Storage, mailSender EmailSender) *Service {
	return &Service{
		log:        logger,
		storage:    storage,
		mailSender: mailSender,
	}
}
