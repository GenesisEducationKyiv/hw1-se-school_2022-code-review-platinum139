package domain

import (
	"subscribers-service/internal/common"
)

type SubscribersServiceImpl struct {
	log             common.Logger
	subscribersRepo Repo
}

func (s SubscribersServiceImpl) Subscribe(subscriber Subscriber) error {
	err := s.subscribersRepo.CreateSubscriber(subscriber)
	if err != nil {
		s.log.Errorf("Unable to create subscriber: %s", subscriber.Email)
		return err
	}
	s.log.Infof("Subscriber %s is created", subscriber.Email)
	return nil
}

func (s SubscribersServiceImpl) GetSubscribers() ([]Subscriber, error) {
	subscribers, err := s.subscribersRepo.GetSubscribers()
	if err != nil {
		s.log.Errorf("Unable to get subscribers from storage:", err)
		return nil, err
	}
	return subscribers, nil
}

func NewSubscribersService(logger common.Logger, subscribersRepo Repo) *SubscribersServiceImpl {
	return &SubscribersServiceImpl{
		log:             logger,
		subscribersRepo: subscribersRepo,
	}
}
