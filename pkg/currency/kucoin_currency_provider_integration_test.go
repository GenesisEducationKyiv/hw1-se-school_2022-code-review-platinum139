package currency

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type KucoinCurrencyProviderIntegrationTestSuite struct {
	suite.Suite
}

func (s *KucoinCurrencyProviderIntegrationTestSuite) TestGetCurrencyRate_Positive() {
	// arrange
	fromCurrency := Btc
	toCurrency := Uah

	service := NewKucoinCurrencyProvider()

	// act
	rate, err := service.GetCurrencyRate(fromCurrency, toCurrency)

	// assert
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), rate)
}

func TestKucoinCurrencyIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(KucoinCurrencyProviderIntegrationTestSuite))
}
