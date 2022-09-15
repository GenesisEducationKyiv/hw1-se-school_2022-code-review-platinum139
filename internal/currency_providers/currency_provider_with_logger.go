package currency_providers

import (
	"bitcoin-service/pkg/currency"
	"log"
)

type ProviderWithLogger struct {
	provider currency.Provider
	logger   *log.Logger
}

func (p *ProviderWithLogger) SetNext(next currency.Provider) {
	p.provider.SetNext(next)
}

func (p *ProviderWithLogger) GetCurrencyRate(from currency.Currency, to currency.Currency) (float64, error) {
	rate, err := p.provider.GetCurrencyRate(from, to)
	if err != nil {
		p.logger.Printf("Failed to get currency rate: %s\n", err)
		return 0, err
	}
	p.logger.Printf("Currency provider response: rate = %f\n", rate)
	return rate, nil
}

func NewCurrencyProviderWithLogger(logger *log.Logger, provider currency.Provider) *ProviderWithLogger {
	return &ProviderWithLogger{
		logger:   logger,
		provider: provider,
	}
}
