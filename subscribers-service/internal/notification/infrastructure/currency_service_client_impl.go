package infrastructure

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"subscribers-service/internal/common"
)

type CurrencyServiceClientImpl struct {
	log                 common.Logger
	currencyServiceHost string
	currencyServicePort string
	rateBitSize         int
}

func (c *CurrencyServiceClientImpl) GetRate() (float64, error) {
	protocol := "http"
	currencySvcUrl := fmt.Sprintf("%s://%s:%s/rate",
		protocol, c.currencyServiceHost, c.currencyServicePort)

	resp, err := http.Get(currencySvcUrl)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	bodyStr := strings.TrimSpace(string(body))
	rate, err := strconv.ParseFloat(bodyStr, c.rateBitSize)
	if err != nil {
		return 0, err
	}

	return rate, nil
}

func NewCurrencyServiceClientImpl(
	log common.Logger,
	currencyServiceHost string,
	currencyServicePort string,
	rateBitSize int,
) *CurrencyServiceClientImpl {
	return &CurrencyServiceClientImpl{
		log:                 log,
		currencyServiceHost: currencyServiceHost,
		currencyServicePort: currencyServicePort,
		rateBitSize:         rateBitSize,
	}
}
