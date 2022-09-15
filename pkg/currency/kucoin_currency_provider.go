package currency

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	kucoinBaseURL   = "https://api.kucoin.com"
	kucoinRateRoute = "api/v1/prices"
)

type KucoinCurrencyProvider struct {
	next Provider
}

func (p *KucoinCurrencyProvider) SetNext(next Provider) {
	p.next = next
}

func (p *KucoinCurrencyProvider) GetCurrencyRate(from Currency, to Currency) (float64, error) {
	rate, err := p.getCurrencyRate(from, to)
	if err != nil && p.next != nil {
		return (p.next).GetCurrencyRate(from, to)
	}
	return rate, err
}

func (p *KucoinCurrencyProvider) getCurrencyRate(from Currency, to Currency) (float64, error) {
	baseQueryParam := fmt.Sprintf("base=%s", to)
	currenciesQueryParam := fmt.Sprintf("currencies=%s", from)
	url := fmt.Sprintf("%s/%s?%s&%s", kucoinBaseURL, kucoinRateRoute, baseQueryParam, currenciesQueryParam)

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
		Code string            `json:"code"`
		Data map[string]string `json:"data"`
	}{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		return 0, err
	}

	key := strings.ToUpper(string(from))
	rateStr, ok := responseBody.Data[key]
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

func NewKucoinCurrencyProvider() *KucoinCurrencyProvider {
	return &KucoinCurrencyProvider{}
}
