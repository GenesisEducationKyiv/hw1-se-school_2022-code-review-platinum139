package currency_providers

import (
	"bitcoin-service/internal/currency/domain"
	"encoding/json"
	"errors"
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

func (p *CoinbaseCurrencyProvider) GetCurrencyRate(from domain.Currency, to domain.Currency) (float64, error) {
	rate, err := p.getCurrencyRate(from, to)
	if err != nil && p.next != nil {
		return (p.next).GetCurrencyRate(from, to)
	}
	return rate, err
}

func (p *CoinbaseCurrencyProvider) getCurrencyRate(from domain.Currency, to domain.Currency) (float64, error) {
	currencyQueryParam := fmt.Sprintf("currency=%s", from)
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

	rateStr, ok := responseBody.Data.Rates[string(to)]
	if !ok {
		return 0, errors.New("rate not found")
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
