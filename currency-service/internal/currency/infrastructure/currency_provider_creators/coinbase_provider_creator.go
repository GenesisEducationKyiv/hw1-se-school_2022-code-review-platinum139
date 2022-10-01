package currency_provider_creators

import (
	"currency-service/internal/currency/domain"
	"currency-service/internal/currency/infrastructure/currency_providers"
)

type CoinbaseProviderCreator struct {
	rateValueBitSize int
}

func (c *CoinbaseProviderCreator) CreateProvider() domain.Provider {
	return currency_providers.NewCoinbaseCurrencyProvider(c.rateValueBitSize)
}

func NewCoinbaseProviderCreator(rateValueBitSize int) *CoinbaseProviderCreator {
	return &CoinbaseProviderCreator{
		rateValueBitSize: rateValueBitSize,
	}
}
