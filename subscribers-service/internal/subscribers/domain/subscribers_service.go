package domain

type SubscribersService interface {
	Subscribe(subscriber Subscriber) error
	GetSubscribers() ([]Subscriber, error)
	Unsubscribe(subscriber Subscriber) error
}
