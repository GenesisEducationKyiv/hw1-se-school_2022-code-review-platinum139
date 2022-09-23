package currency_provider_creators

import (
	"bitcoin-service/internal/currency/domain"
	"bitcoin-service/internal/currency/infrastructure/currency_providers"
	"log"
	"time"
)

type CoingateProviderCreator struct {
	log              *log.Logger
	cachingPeriodMin int
	rateValueBitSize int
}

func (c *CoingateProviderCreator) CreateProvider() domain.Provider {
	provider1 := domain.NewCachedCurrencyProvider(
		time.Duration(c.cachingPeriodMin)*time.Minute,
		c.rateValueBitSize,
		domain.NewCurrencyProviderWithLogger(c.log, currency_providers.NewCoingateCurrencyProvider(c.rateValueBitSize)))

	provider2 := domain.NewCachedCurrencyProvider(
		time.Duration(c.cachingPeriodMin)*time.Minute,
		c.rateValueBitSize,
		domain.NewCurrencyProviderWithLogger(c.log, currency_providers.NewCoinbaseCurrencyProvider(c.rateValueBitSize)))

	provider3 := domain.NewCachedCurrencyProvider(
		time.Duration(c.cachingPeriodMin)*time.Minute,
		c.rateValueBitSize,
		domain.NewCurrencyProviderWithLogger(c.log, currency_providers.NewKucoinCurrencyProvider(c.rateValueBitSize)))

	provider1.SetNext(provider2)
	provider2.SetNext(provider3)

	return provider1
}

func NewCoingateProviderCreator(log *log.Logger, cachingPeriodMin, rateValueBitSize int) *CoingateProviderCreator {
	return &CoingateProviderCreator{
		log:              log,
		cachingPeriodMin: cachingPeriodMin,
		rateValueBitSize: rateValueBitSize,
	}
}
