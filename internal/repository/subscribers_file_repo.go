package repository

import (
	"bitcoin-service/internal/subscribers"
)

type Storage interface {
	Add(record string) error
	GetAll() ([]string, error)
}

type SubscribersFileRepo struct {
	storage Storage
}

func (repo *SubscribersFileRepo) CreateSubscriber(subscriber subscribers.Subscriber) error {
	return repo.storage.Add(subscriber.Email)
}

func (repo *SubscribersFileRepo) GetSubscribers() ([]subscribers.Subscriber, error) {
	emails, err := repo.storage.GetAll()
	if err != nil {
		return nil, err
	}

	var resultSubscribers []subscribers.Subscriber
	for _, email := range emails {
		resultSubscribers = append(resultSubscribers, subscribers.Subscriber{
			Email: email,
		})
	}

	return resultSubscribers, nil
}

func NewSubscribersFileRepo(storage Storage) *SubscribersFileRepo {
	return &SubscribersFileRepo{
		storage: storage,
	}
}
