package controller

import subscribers "subscribers-service/internal/subscribers/domain"

func SubscriberToDTO(subscriber subscribers.Subscriber) SubscriberDTO {
	return SubscriberDTO{
		Email: subscriber.Email,
	}
}

func SubscriberFromDTO(subscriberDTO SubscriberDTO) subscribers.Subscriber {
	return subscribers.Subscriber{
		Email: subscriberDTO.Email,
	}
}
