package domain

import (
	"fmt"
	"subscribers-service/config"
	"subscribers-service/internal/common"
	subscribersDomain "subscribers-service/internal/subscribers/domain"
)

type ServiceImpl struct {
	log                common.Logger
	cfg                *config.AppConfig
	messageSender      MessageSender
	currencySvcClient  CurrencyServiceClient
	subscribersService subscribersDomain.SubscribersService
}

func (s *ServiceImpl) Notify() error {
	fromCurrency := Currency(s.cfg.FromCurrency)
	toCurrency := Currency(s.cfg.ToCurrency)

	rate, err := s.currencySvcClient.GetRate()
	if err != nil {
		s.log.Errorf("Failed to get rate from currency service: %s", err)
		return err
	}
	s.log.Infof("Got rate from currency service: %f", rate)

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
			s.log.Errorf("Unable to send mails via mail service for %s: %s", subscriber, err)
			failedSubscribers = append(failedSubscribers, subscriber.Email)
		}
	}

	if len(failedSubscribers) != 0 {
		return SendMessageError{FailedSubscribers: failedSubscribers}
	}

	return nil
}

func NewNotificationService(
	log common.Logger,
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
