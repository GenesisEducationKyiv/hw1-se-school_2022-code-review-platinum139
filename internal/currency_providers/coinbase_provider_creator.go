package currency_providers

import (
	"bitcoin-service/pkg/currency"
	"log"
	"time"
)

type CoinbaseProviderCreator struct {
	log              *log.Logger
	cachingPeriodMin int
}

func (c *CoinbaseProviderCreator) CreateProvider() currency.Provider {
	provider1 := NewCachedCurrencyProvider(
		time.Duration(c.cachingPeriodMin)*time.Minute,
		NewCurrencyProviderWithLogger(c.log, currency.NewCoinbaseCurrencyProvider()))

	provider2 := NewCachedCurrencyProvider(
		time.Duration(c.cachingPeriodMin)*time.Minute,
		NewCurrencyProviderWithLogger(c.log, currency.NewCoingateCurrencyProvider()))

	provider3 := NewCachedCurrencyProvider(
		time.Duration(c.cachingPeriodMin)*time.Minute,
		NewCurrencyProviderWithLogger(c.log, currency.NewKucoinCurrencyProvider()))

	provider1.SetNext(provider2)
	provider2.SetNext(provider3)

	return provider1
}

func NewCoinbaseProviderCreator(log *log.Logger, cachingPeriodMin int) *CoinbaseProviderCreator {
	return &CoinbaseProviderCreator{
		log:              log,
		cachingPeriodMin: cachingPeriodMin,
	}
}
