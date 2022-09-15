package currency

import (
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
	next Provider
}

func (p *CoinbaseCurrencyProvider) SetNext(next Provider) {
	p.next = next
}

func (p *CoinbaseCurrencyProvider) GetCurrencyRate(from Currency, to Currency) (float64, error) {
	rate, err := p.getCurrencyRate(from, to)
	if err != nil && p.next != nil {
		return (p.next).GetCurrencyRate(from, to)
	}
	return rate, err
}

func (p *CoinbaseCurrencyProvider) getCurrencyRate(from Currency, to Currency) (float64, error) {
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

	bitSize := 64
	rate, err := strconv.ParseFloat(rateStr, bitSize)
	if err != nil {
		return 0, err
	}

	return rate, nil
}

func NewCoinbaseCurrencyProvider() *CoinbaseCurrencyProvider {
	return &CoinbaseCurrencyProvider{}
}
