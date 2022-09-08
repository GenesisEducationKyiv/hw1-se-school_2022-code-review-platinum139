package currency

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CoingateCurrencyProviderServiceIntegrationTestSuite struct {
	suite.Suite
}

func (s *CoingateCurrencyProviderServiceIntegrationTestSuite) TestGetCurrencyRate_Positive() {
	// arrange
	fromCurrency := Btc
	toCurrency := Uah

	service := NewCurrencyService()

	// act
	rate, err := service.GetCurrencyRate(fromCurrency, toCurrency)

	// assert
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), rate)
}

func TestCurrencyIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(CoingateCurrencyProviderServiceIntegrationTestSuite))
}
