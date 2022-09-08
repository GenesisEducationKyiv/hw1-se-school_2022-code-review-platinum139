package currency

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const (
	baseURL   = "https://api.coingate.com"
	rateRoute = "v2/rates/merchant"
)

type CoingateCurrencyProvider struct {
}

func (s *CoingateCurrencyProvider) GetCurrencyRate(from Currency, to Currency) (float64, error) {
	url := fmt.Sprintf("%s/%s/%s/%s", baseURL, rateRoute, from, to)

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

func NewCurrencyService() *CoingateCurrencyProvider {
	return &CoingateCurrencyProvider{}
}
