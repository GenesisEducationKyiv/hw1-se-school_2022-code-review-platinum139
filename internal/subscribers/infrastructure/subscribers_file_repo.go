package infrastructure

import (
	"bitcoin-service/internal/subscribers/domain"
	"bitcoin-service/pkg/file_storage"
)

type SubscribersFileRepo struct {
	storage file_storage.FileStorage
}

func (repo *SubscribersFileRepo) CreateSubscriber(subscriber domain.Subscriber) error {
	return repo.storage.Add(subscriber.Email)
}

func (repo *SubscribersFileRepo) GetSubscribers() ([]domain.Subscriber, error) {
	emails, err := repo.storage.GetAll()
	if err != nil {
		return nil, err
	}

	var resultSubscribers []domain.Subscriber
	for _, email := range emails {
		resultSubscribers = append(resultSubscribers, domain.Subscriber{
			Email: email,
		})
	}

	return resultSubscribers, nil
}

func NewSubscribersFileRepo(storage file_storage.FileStorage) *SubscribersFileRepo {
	return &SubscribersFileRepo{
		storage: storage,
	}
}
