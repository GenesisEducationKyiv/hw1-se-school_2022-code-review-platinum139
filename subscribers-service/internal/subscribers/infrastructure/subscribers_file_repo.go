package infrastructure

import (
	"subscribers-service/internal/subscribers/domain"
	"subscribers-service/pkg/file_storage"
)

type SubscribersFileRepo struct {
	storage file_storage.FileStorage
}

func (repo *SubscribersFileRepo) CreateSubscriber(subscriber domain.Subscriber) error {
	return repo.storage.Add(subscriber.TransactionID, subscriber.Email)
}

func (repo *SubscribersFileRepo) GetSubscribers() ([]domain.Subscriber, error) {
	records, err := repo.storage.GetAll()
	if err != nil {
		return nil, err
	}

	resultSubscribers := make([]domain.Subscriber, len(records))
	for i, elems := range records {
		resultSubscribers[i] = domain.Subscriber{
			TransactionID: elems[0],
			Email:         elems[1],
		}
	}

	return resultSubscribers, nil
}

func (repo *SubscribersFileRepo) DeleteSubscriber(subscriber domain.Subscriber) error {
	return repo.storage.Delete(subscriber.TransactionID)
}

func NewSubscribersFileRepo(storage file_storage.FileStorage) *SubscribersFileRepo {
	return &SubscribersFileRepo{
		storage: storage,
	}
}
