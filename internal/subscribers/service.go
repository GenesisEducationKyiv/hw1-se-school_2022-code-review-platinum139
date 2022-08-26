package subscribers

import (
	"bitcoin-service/pkg/emails"
	"bitcoin-service/pkg/rates"
	"fmt"
	"log"
)

const (
	uah = "UAH"
	btc = "BTC"
)

type Storage interface {
	Add(record string) error
	GetAll() ([]string, error)
}

type Service struct {
	log         *log.Logger
	storage     Storage
	mailService *emails.Service
}

func (s Service) Add(subscriber string) error {
	return s.storage.Add(subscriber)
}

func (s Service) SendEmails() error {
	subscribers, err := s.getAll()
	if err != nil {
		s.log.Print("Unable to get subscribers from storage:", err)
		return err
	}

	rate, err := rates.GetCurrencyRate(btc, uah)
	if err != nil {
		s.log.Print("Unable to get btc-to-uah rate:", err)
		return err
	}

	var failedEmails []string
	for _, subscriber := range subscribers {
		message := fmt.Sprintf("Rate = %f", rate)

		err := s.mailService.SendEmail(subscriber, "BTC to UAH rate", message)
		if err != nil {
			s.log.Print(fmt.Sprintf(
				"Unable to send mails via mail service for %s: %s", subscriber, err))
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

func NewService(logger *log.Logger, storage Storage, mailService *emails.Service) *Service {
	return &Service{
		log:         logger,
		storage:     storage,
		mailService: mailService,
	}
}
