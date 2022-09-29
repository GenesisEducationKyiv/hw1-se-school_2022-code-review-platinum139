package currency_provider_creators

import (
	"currency-service/internal/currency/domain"
	"currency-service/internal/currency/infrastructure/currency_providers"
)

type CoingateProviderCreator struct {
	rateValueBitSize int
}

func (c *CoingateProviderCreator) CreateProvider() domain.Provider {
	return currency_providers.NewCoingateCurrencyProvider(c.rateValueBitSize)
}

func NewCoingateProviderCreator(rateValueBitSize int) *CoingateProviderCreator {
	return &CoingateProviderCreator{
		rateValueBitSize: rateValueBitSize,
	}
}
