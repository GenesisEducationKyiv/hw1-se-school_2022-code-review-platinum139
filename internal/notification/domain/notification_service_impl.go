package domain

import (
	"bitcoin-service/config"
	currencyDomain "bitcoin-service/internal/currency/domain"
	subscribersDomain "bitcoin-service/internal/subscribers/domain"
	"fmt"
	"log"
)

type ServiceImpl struct {
	log                *log.Logger
	cfg                *config.AppConfig
	messageSender      MessageSender
	currencyService    currencyDomain.CurrencyService
	subscribersService subscribersDomain.SubscribersService
}

func (s *ServiceImpl) Notify() error {
	fromCurrency := currencyDomain.Currency(s.cfg.FromCurrency)
	toCurrency := currencyDomain.Currency(s.cfg.ToCurrency)

	rate, err := s.currencyService.GetCurrencyRate(fromCurrency, toCurrency)
	if err != nil {
		return err
	}

	message := Message{
		Subject: fmt.Sprintf("%s to %s rate", fromCurrency, toCurrency),
		Text:    fmt.Sprintf("Rate = %.2f", rate),
	}

	subscribers, err := s.subscribersService.GetSubscribers()
	if err != nil {
		return err
	}

	var failedSubscribers []string
	for _, subscriber := range subscribers {
		err := s.messageSender.Send(subscriber.Email, message.Subject, message.Text)
		if err != nil {
			s.log.Printf("Unable to send mails via mail service for %s: %s", subscriber, err)
			failedSubscribers = append(failedSubscribers, subscriber.Email)
		}
	}

	if len(failedSubscribers) != 0 {
		return SendMessageError{FailedSubscribers: failedSubscribers}
	}

	return nil
}

func NewNotificationService(
	log *log.Logger,
	cfg *config.AppConfig,
	messageSender MessageSender,
	currencyService currencyDomain.CurrencyService,
	subscribersService subscribersDomain.SubscribersService,
) *ServiceImpl {
	return &ServiceImpl{
		log:                log,
		cfg:                cfg,
		messageSender:      messageSender,
		currencyService:    currencyService,
		subscribersService: subscribersService,
	}
}
