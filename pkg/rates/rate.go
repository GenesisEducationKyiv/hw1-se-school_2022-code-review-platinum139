package rates

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func GetCurrencyRate(from string, to string) (float64, error) {
	url := fmt.Sprintf("https://api.coingate.com/v2/rates/merchant/%s/%s", from, to)

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

	rate, err := strconv.ParseFloat(bodyString, 64)
	if err != nil {
		return 0, err
	}

	return rate, nil
}
