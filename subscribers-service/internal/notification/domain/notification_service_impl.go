package domain

import (
	"fmt"
	"log"
	"subscribers-service/config"
	"subscribers-service/internal/notification/common"
	subscribersDomain "subscribers-service/internal/subscribers/domain"
)

type ServiceImpl struct {
	log                *log.Logger
	cfg                *config.AppConfig
	messageSender      MessageSender
	currencySvcClient  CurrencyServiceClient
	subscribersService subscribersDomain.SubscribersService
}

func (s *ServiceImpl) Notify() error {
	fromCurrency := common.Currency(s.cfg.FromCurrency)
	toCurrency := common.Currency(s.cfg.ToCurrency)

	rate, err := s.currencySvcClient.GetRate()
	if err != nil {
		s.log.Printf("Failed to get rate from currency service: %s", err)
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
	currencySvcClient CurrencyServiceClient,
	subscribersService subscribersDomain.SubscribersService,
) *ServiceImpl {
	return &ServiceImpl{
		log:                log,
		cfg:                cfg,
		messageSender:      messageSender,
		currencySvcClient:  currencySvcClient,
		subscribersService: subscribersService,
	}
}
