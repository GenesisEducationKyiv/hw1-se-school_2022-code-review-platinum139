package domain

type Repo interface {
	CreateSubscriber(subscriber Subscriber) error
	GetSubscribers() ([]Subscriber, error)
	DeleteSubscriber(subscriber Subscriber) error
}
