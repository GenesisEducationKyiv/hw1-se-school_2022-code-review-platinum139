package currency_provider_creators

import (
	"currency-service/internal/currency/domain"
	"currency-service/internal/currency/infrastructure/currency_providers"
)

type KucoinProviderCreator struct {
	rateValueBitSize int
}

func (c *KucoinProviderCreator) CreateProvider() domain.Provider {
	return currency_providers.NewKucoinCurrencyProvider(c.rateValueBitSize)
}

func NewKucoinProviderCreator(rateValueBitSize int) *KucoinProviderCreator {
	return &KucoinProviderCreator{
		rateValueBitSize: rateValueBitSize,
	}
}
