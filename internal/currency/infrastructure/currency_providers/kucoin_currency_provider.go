package currency_providers

import (
	"bitcoin-service/internal/currency/domain"
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
	next             domain.Provider
	rateValueBitSize int
}

func (p *KucoinCurrencyProvider) SetNext(next domain.Provider) {
	p.next = next
}

func (p *KucoinCurrencyProvider) GetCurrencyRate(from domain.Currency, to domain.Currency) (float64, error) {
	rate, err := p.getCurrencyRate(from, to)
	if err != nil && p.next != nil {
		return (p.next).GetCurrencyRate(from, to)
	}
	return rate, err
}

func (p *KucoinCurrencyProvider) getCurrencyRate(from domain.Currency, to domain.Currency) (float64, error) {
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

	rate, err := strconv.ParseFloat(rateStr, p.rateValueBitSize)
	if err != nil {
		return 0, err
	}

	return rate, nil
}

func NewKucoinCurrencyProvider(rateValueBitSize int) *KucoinCurrencyProvider {
	return &KucoinCurrencyProvider{
		rateValueBitSize: rateValueBitSize,
	}
}
