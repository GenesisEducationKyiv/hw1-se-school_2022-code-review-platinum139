package domain

import (
	"fmt"
	"strconv"
	"time"
)

type CachedValue struct {
	Value     string
	ExpiresAt time.Time
}

type CachedCurrencyProvider struct {
	provider         Provider
	cache            map[string]CachedValue
	cachingPeriod    time.Duration
	rateValueBitSize int
}

func (p *CachedCurrencyProvider) SetNext(next Provider) {
	p.provider.SetNext(next)
}

func (p *CachedCurrencyProvider) GetCurrencyRate(from Currency, to Currency) (float64, error) {
	key := fmt.Sprintf("%s:%s", from, to)

	cachedValue, ok := p.cache[key]
	if ok && cachedValue.ExpiresAt.After(time.Now()) {
		rate, err := strconv.ParseFloat(cachedValue.Value, p.rateValueBitSize)
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

func NewCachedCurrencyProvider(
	cachingPeriod time.Duration,
	rateValueBitSize int,
	provider Provider,
) *CachedCurrencyProvider {
	return &CachedCurrencyProvider{
		provider:         provider,
		cache:            make(map[string]CachedValue),
		cachingPeriod:    cachingPeriod,
		rateValueBitSize: rateValueBitSize,
	}
}
