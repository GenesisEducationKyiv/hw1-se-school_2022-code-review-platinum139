package rates

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const baseURL = "https://api.coingate.com/v2/rates/merchant"

func GetCurrencyRate(from string, to string) (float64, error) {
	url := fmt.Sprintf("%s/%s/%s", baseURL, from, to)

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
