package currency_providers

import (
	"currency-service/internal/currency/domain"
	"currency-service/internal/currency/infrastructure/currency_provider_errors"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const (
	coinbaseBaseURL   = "https://api.coinbase.com"
	coinbaseRateRoute = "v2/exchange-rates"
)

type CoinbaseCurrencyProvider struct {
	next             domain.Provider
	rateValueBitSize int
}

func (p *CoinbaseCurrencyProvider) SetNext(next domain.Provider) {
	p.next = next
}

func (p *CoinbaseCurrencyProvider) GetCurrencyRate(
	fromCurrency domain.Currency,
	toCurrency domain.Currency,
) (float64, error) {
	rate, err := p.getCurrencyRate(fromCurrency, toCurrency)
	if err != nil && p.next != nil {
		return (p.next).GetCurrencyRate(fromCurrency, toCurrency)
	}
	return rate, err
}

func (p *CoinbaseCurrencyProvider) getCurrencyRate(
	fromCurrency domain.Currency,
	toCurrency domain.Currency,
) (float64, error) {
	currencyQueryParam := fmt.Sprintf("currency=%s", fromCurrency)
	url := fmt.Sprintf("%s/%s?%s", coinbaseBaseURL, coinbaseRateRoute, currencyQueryParam)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	responseBody := struct {
		Data struct {
			Currency string
			Rates    map[string]string
		}
	}{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		return 0, err
	}

	rateStr, ok := responseBody.Data.Rates[string(toCurrency)]
	if !ok {
		return 0, currency_provider_errors.RateNotFoundError{}
	}

	rate, err := strconv.ParseFloat(rateStr, p.rateValueBitSize)
	if err != nil {
		return 0, err
	}

	return rate, nil
}

func NewCoinbaseCurrencyProvider(rateValueBitSize int) *CoinbaseCurrencyProvider {
	return &CoinbaseCurrencyProvider{
		rateValueBitSize: rateValueBitSize,
	}
}
