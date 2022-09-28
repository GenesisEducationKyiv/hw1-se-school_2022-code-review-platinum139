package currency_providers

import (
	"bitcoin-service/internal/currency/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CoinbaseCurrencyProviderIntegrationTestSuite struct {
	suite.Suite
}

func (s *CoinbaseCurrencyProviderIntegrationTestSuite) TestGetCurrencyRate_Positive() {
	// arrange
	fromCurrency := domain.Btc
	toCurrency := domain.Uah

	bitSize := 64
	service := NewCoinbaseCurrencyProvider(bitSize)

	// act
	rate, err := service.GetCurrencyRate(fromCurrency, toCurrency)

	// assert
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), rate)
}

func TestCoinbaseCurrencyIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(CoinbaseCurrencyProviderIntegrationTestSuite))
}
