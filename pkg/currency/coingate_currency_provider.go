package currency

import (
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
	next Provider
}

func (p *CoingateCurrencyProvider) SetNext(next Provider) {
	p.next = next
}

func (p *CoingateCurrencyProvider) GetCurrencyRate(from Currency, to Currency) (float64, error) {
	rate, err := p.getCurrencyRate(from, to)
	if err != nil && p.next != nil {
		return (p.next).GetCurrencyRate(from, to)
	}
	return rate, err
}

func (s *CoingateCurrencyProvider) getCurrencyRate(from Currency, to Currency) (float64, error) {
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

	bitSize := 64
	rate, err := strconv.ParseFloat(bodyString, bitSize)
	if err != nil {
		return 0, err
	}

	return rate, nil
}

func NewCoingateCurrencyProvider() *CoingateCurrencyProvider {
	return &CoingateCurrencyProvider{}
}
