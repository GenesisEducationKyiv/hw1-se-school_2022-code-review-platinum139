package domain

import (
	"log"
)

type ProviderWithLogger struct {
	provider Provider
	logger   *log.Logger
}

func (p *ProviderWithLogger) SetNext(next Provider) {
	p.provider.SetNext(next)
}

func (p *ProviderWithLogger) GetCurrencyRate(from Currency, to Currency) (float64, error) {
	rate, err := p.provider.GetCurrencyRate(from, to)
	if err != nil {
		p.logger.Printf("Failed to get currency rate: %s\n", err)
		return 0, err
	}
	p.logger.Printf("Currency provider response: rate = %f\n", rate)
	return rate, nil
}

func NewCurrencyProviderWithLogger(logger *log.Logger, provider Provider) *ProviderWithLogger {
	return &ProviderWithLogger{
		logger:   logger,
		provider: provider,
	}
}
