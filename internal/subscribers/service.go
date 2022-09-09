package subscribers

import (
	"bitcoin-service/internal/common"
	"log"
)

type Repo interface {
	CreateSubscriber(subscriber Subscriber) error
	GetSubscribers() ([]Subscriber, error)
}

type MessageSender interface {
	Send(receiver string, subject string, text string) error
}

type Service struct {
	log             *log.Logger
	subscribersRepo Repo
	messageSender   MessageSender
}

func (s Service) Subscribe(subscriber Subscriber) error {
	return s.subscribersRepo.CreateSubscriber(subscriber)
}

func (s Service) SendEmails(message common.Message) error {
	subscribers, err := s.subscribersRepo.GetSubscribers()
	if err != nil {
		s.log.Print("Unable to get subscribers from storage:", err)
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

func NewSubscribersService(logger *log.Logger, subscribersRepo Repo, mailSender MessageSender) *Service {
	return &Service{
		log:             logger,
		subscribersRepo: subscribersRepo,
		messageSender:   mailSender,
	}
}
