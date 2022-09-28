package currency_providers

import (
	"bitcoin-service/internal/currency/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CoingateCurrencyProviderIntegrationTestSuite struct {
	suite.Suite
}

func (s *CoingateCurrencyProviderIntegrationTestSuite) TestGetCurrencyRate_Positive() {
	// arrange
	fromCurrency := domain.Btc
	toCurrency := domain.Uah

	bitSize := 64
	service := NewCoingateCurrencyProvider(bitSize)

	// act
	rate, err := service.GetCurrencyRate(fromCurrency, toCurrency)

	// assert
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), rate)
}

func TestCoingateCurrencyIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(CoingateCurrencyProviderIntegrationTestSuite))
}
