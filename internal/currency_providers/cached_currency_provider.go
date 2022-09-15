package currency_providers

import (
	"bitcoin-service/pkg/currency"
	"fmt"
	"strconv"
	"time"
)

type CachedValue struct {
	Value     string
	ExpiresAt time.Time
}

type CachedCurrencyProvider struct {
	provider      currency.Provider
	cache         map[string]CachedValue
	cachingPeriod time.Duration
}

func (p *CachedCurrencyProvider) SetNext(next currency.Provider) {
	p.provider.SetNext(next)
}

func (p *CachedCurrencyProvider) GetCurrencyRate(from currency.Currency, to currency.Currency) (float64, error) {
	key := fmt.Sprintf("%s:%s", from, to)

	cachedValue, ok := p.cache[key]
	if ok && cachedValue.ExpiresAt.After(time.Now()) {
		bitSize := 64
		rate, err := strconv.ParseFloat(cachedValue.Value, bitSize)
		if err != nil {
			return 0, err
		}
		return rate, nil
	}

	rate, err := p.provider.GetCurrencyRate(from, to)
	if err != nil {
		return 0, err
	}

	p.cache[key] = CachedValue{
		Value:     fmt.Sprintf("%f", rate),
		ExpiresAt: time.Now().Add(p.cachingPeriod),
	}
	return rate, nil
}

func NewCachedCurrencyProvider(cachingPeriod time.Duration, provider currency.Provider) *CachedCurrencyProvider {
	return &CachedCurrencyProvider{
		provider:      provider,
		cache:         make(map[string]CachedValue),
		cachingPeriod: cachingPeriod,
	}
}
