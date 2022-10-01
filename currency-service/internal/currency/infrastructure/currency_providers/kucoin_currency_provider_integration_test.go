package currency_providers

import (
	"currency-service/internal/currency/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type KucoinCurrencyProviderIntegrationTestSuite struct {
	suite.Suite
}

func (s *KucoinCurrencyProviderIntegrationTestSuite) TestGetCurrencyRate_Positive() {
	// arrange
	fromCurrency := domain.Btc
	toCurrency := domain.Uah

	bitSize := 64
	service := NewKucoinCurrencyProvider(bitSize)

	// act
	rate, err := service.GetCurrencyRate(fromCurrency, toCurrency)

	// assert
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), rate)
}

func TestKucoinCurrencyIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(KucoinCurrencyProviderIntegrationTestSuite))
}
