package currency

import (
	"bitcoin-service/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CurrencyServiceIntegrationTestSuite struct {
	suite.Suite
}

func (s *CurrencyServiceIntegrationTestSuite) TestGetCurrencyRate_Positive() {
	// arrange
	appConfig, err := config.NewAppConfig(".env.test")
	assert.NoError(s.T(), err)

	fromCurrency := Btc
	toCurrency := Uah

	service := NewCurrencyService(appConfig.CurrencyServiceBaseURL)

	// act
	rate, err := service.GetCurrencyRate(fromCurrency, toCurrency)

	// assert
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), rate)
}

func TestCurrencyIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyServiceIntegrationTestSuite))
}
