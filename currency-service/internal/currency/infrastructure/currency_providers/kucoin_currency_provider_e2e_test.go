package currency_providers

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

func TestKucoinE2E(t *testing.T) {
	url := "https://api.kucoin.com/api/v1/prices?base=UAH&currencies=BTC"

	expectedBody := `
	{
		"code": "[\d]+",
		"data": {
			"BTC": "[\d]+\.[\d]+"
		}
	}`
	expectedBody = strings.Join(strings.Fields(expectedBody), "")
	t.Log(expectedBody)

	re := regexp.MustCompile(expectedBody)

	resp, err := http.Get(url)
	assert.NoError(t, err)

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	t.Log(string(bodyBytes))

	matched := re.MatchString(string(bodyBytes))
	assert.EqualValues(t, true, matched)
}
