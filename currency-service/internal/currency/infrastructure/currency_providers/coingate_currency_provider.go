package currency_providers

import (
	"currency-service/internal/currency/domain"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const (
	coingateBaseURL   = "https://api.coingate.com"
	coingateRateRoute = "v2/rates/merchant"
)

type CoingateCurrencyProvider struct {
	next             domain.Provider
	rateValueBitSize int
}

func (p *CoingateCurrencyProvider) SetNext(next domain.Provider) {
	p.next = next
}

func (p *CoingateCurrencyProvider) GetCurrencyRate(from domain.Currency, to domain.Currency) (float64, error) {
	rate, err := p.getCurrencyRate(from, to)
	if err != nil && p.next != nil {
		return (p.next).GetCurrencyRate(from, to)
	}
	return rate, err
}

func (s *CoingateCurrencyProvider) getCurrencyRate(from domain.Currency, to domain.Currency) (float64, error) {
	url := fmt.Sprintf("%s/%s/%s/%s", coingateBaseURL, coingateRateRoute, from, to)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	bodyString := string(bodyBytes)

	rate, err := strconv.ParseFloat(bodyString, s.rateValueBitSize)
	if err != nil {
		return 0, err
	}

	return rate, nil
}

func NewCoingateCurrencyProvider(rateValueBitSize int) *CoingateCurrencyProvider {
	return &CoingateCurrencyProvider{
		rateValueBitSize: rateValueBitSize,
	}
}
