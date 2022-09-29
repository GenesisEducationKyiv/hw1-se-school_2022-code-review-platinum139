package domain

import (
	"log"
)

type SubscribersServiceImpl struct {
	log             *log.Logger
	subscribersRepo Repo
}

func (s SubscribersServiceImpl) Subscribe(subscriber Subscriber) error {
	return s.subscribersRepo.CreateSubscriber(subscriber)
}

func (s SubscribersServiceImpl) GetSubscribers() ([]Subscriber, error) {
	subscribers, err := s.subscribersRepo.GetSubscribers()
	if err != nil {
		s.log.Print("Unable to get subscribers from storage:", err)
		return nil, err
	}
	return subscribers, nil
}

func NewSubscribersService(logger *log.Logger, subscribersRepo Repo) *SubscribersServiceImpl {
	return &SubscribersServiceImpl{
		log:             logger,
		subscribersRepo: subscribersRepo,
	}
}
