package controller

import (
	"currency-service/config"
	"currency-service/internal/currency/domain"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CurrencyProviderMock struct {
	mock.Mock
}

func (s *CurrencyProviderMock) GetCurrencyRate(from domain.Currency, to domain.Currency) (float64, error) {
	args := s.Called(from, to)
	return args.Get(0).(float64), args.Error(1)
}

func (s *CurrencyProviderMock) SetNext(next domain.Provider) {
	return
}

type CurrencyControllerUnitTestSuite struct {
	suite.Suite
	config *config.AppConfig
	logger *log.Logger
}

func (s *CurrencyControllerUnitTestSuite) SetupSuite() {
	appConfig, err := config.NewAppConfig(".env.test")
	assert.NoError(s.T(), err)
	s.config = appConfig

	s.logger = log.New(os.Stdout, "", appConfig.LogLevel)
}

func (s *CurrencyControllerUnitTestSuite) TestGetRateHandler_StatusOK() {
	// arrange
	currencyProvider := new(CurrencyProviderMock)
	currencyProvider.On("GetCurrencyRate", domain.Btc, domain.Uah).Return(50000.0, nil)

	currencyService := domain.NewCurrencyService(currencyProvider)
	currencyController := NewCurrencyController(s.logger, currencyService)

	request := httptest.NewRequest(http.MethodGet, "/rate", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := currencyController.GetRateHandler(ctx)

	// assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, responseRecorder.Code)
}

func (s *CurrencyControllerUnitTestSuite) TestGetRateHandler_StatusBadRequest() {
	// arrange
	expectedBody := map[string]string{
		"error": "invalid status value",
	}

	currencyProvider := new(CurrencyProviderMock)
	currencyProvider.On("GetCurrencyRate", domain.Btc, domain.Uah).
		Return(0.0, http.ErrServerClosed)

	currencyService := domain.NewCurrencyService(currencyProvider)
	currencyController := NewCurrencyController(s.logger, currencyService)

	request := httptest.NewRequest(http.MethodGet, "/rate", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := currencyController.GetRateHandler(ctx)
	httpErr, ok := err.(*echo.HTTPError)

	// assert
	assert.Equal(s.T(), true, ok)
	assert.Equal(s.T(), expectedBody, httpErr.Message)
	assert.Equal(s.T(), http.StatusBadRequest, httpErr.Code)
}

func TestCurrencyControllerUnitTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyControllerUnitTestSuite))
}
