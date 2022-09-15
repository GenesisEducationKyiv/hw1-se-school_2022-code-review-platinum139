package currency_providers

import (
	"bitcoin-service/pkg/currency"
	"log"
	"time"
)

type CoingateProviderCreator struct {
	log              *log.Logger
	cachingPeriodMin int
}

func (c *CoingateProviderCreator) CreateProvider() currency.Provider {
	provider1 := NewCachedCurrencyProvider(
		time.Duration(c.cachingPeriodMin)*time.Minute,
		NewCurrencyProviderWithLogger(c.log, currency.NewCoingateCurrencyProvider()))

	provider2 := NewCachedCurrencyProvider(
		time.Duration(c.cachingPeriodMin)*time.Minute,
		NewCurrencyProviderWithLogger(c.log, currency.NewCoinbaseCurrencyProvider()))

	provider3 := NewCachedCurrencyProvider(
		time.Duration(c.cachingPeriodMin)*time.Minute,
		NewCurrencyProviderWithLogger(c.log, currency.NewKucoinCurrencyProvider()))

	provider1.SetNext(provider2)
	provider2.SetNext(provider3)

	return provider1
}

func NewCoingateProviderCreator(log *log.Logger, cachingPeriodMin int) *CoingateProviderCreator {
	return &CoingateProviderCreator{
		log:              log,
		cachingPeriodMin: cachingPeriodMin,
	}
}
