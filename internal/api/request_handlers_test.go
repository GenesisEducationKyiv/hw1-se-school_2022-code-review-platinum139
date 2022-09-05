package api

import (
	"bitcoin-service/config"
	"bitcoin-service/internal/subscribers"
	"bitcoin-service/pkg/currency"
	"bitcoin-service/pkg/storage"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SubscribersServiceMock struct {
	mock.Mock
}

func (s *SubscribersServiceMock) Add(subscriber string) error {
	args := s.Called(subscriber)
	return args.Error(0)
}

func (s *SubscribersServiceMock) SendEmails(rate float64, fromCurrency, toCurrency string) error {
	args := s.Called(rate, fromCurrency, toCurrency)
	return args.Error(0)
}

type CurrencyServiceMock struct {
	mock.Mock
}

func (s *CurrencyServiceMock) GetCurrencyRate(from currency.Currency, to currency.Currency) (float64, error) {
	args := s.Called(from, to)
	return args.Get(0).(float64), args.Error(1)
}

type RequestHandlersUnitTestSuite struct {
	suite.Suite
}

func (s *RequestHandlersUnitTestSuite) TestGetRateHandler_StatusOK() {
	// arrange
	appConfig, err := config.NewAppConfig(".env.test")
	assert.NoError(s.T(), err)

	logger := log.New(os.Stdout, "", appConfig.LogLevel)
	subscribersService := new(SubscribersServiceMock)

	currencyService := new(CurrencyServiceMock)
	currencyService.On("GetCurrencyRate", currency.Btc, currency.Uah).Return(50000.0, nil)

	server := NewServer(logger, appConfig, subscribersService, currencyService)

	request := httptest.NewRequest(http.MethodGet, "/rate", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err = server.GetRateHandler(ctx)

	// assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, responseRecorder.Code)
}

func (s *RequestHandlersUnitTestSuite) TestGetRateHandler_StatusBadRequest() {
	// arrange
	expectedBody := map[string]string{
		"error": "invalid status value",
	}

	appConfig, err := config.NewAppConfig(".env.test")
	assert.NoError(s.T(), err)

	logger := log.New(os.Stdout, "", appConfig.LogLevel)
	subscribersService := new(SubscribersServiceMock)

	currencyService := new(CurrencyServiceMock)
	currencyService.On("GetCurrencyRate", currency.Btc, currency.Uah).Return(0.0, http.ErrServerClosed)

	server := NewServer(logger, appConfig, subscribersService, currencyService)

	request := httptest.NewRequest(http.MethodGet, "/rate", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err = server.GetRateHandler(ctx)
	httpErr, ok := err.(*echo.HTTPError)

	// assert
	assert.Equal(s.T(), true, ok)
	assert.Equal(s.T(), expectedBody, httpErr.Message)
	assert.Equal(s.T(), http.StatusBadRequest, httpErr.Code)
}

func (s *RequestHandlersUnitTestSuite) TestSubscribeHandler_EmailNotProvided() {
	// arrange
	expectedBody := map[string]string{
		"error": "email must be provided",
	}
	appConfig, err := config.NewAppConfig(".env.test")
	assert.NoError(s.T(), err)

	logger := log.New(os.Stdout, "", appConfig.LogLevel)
	subscribersService := new(SubscribersServiceMock)
	currencyService := new(CurrencyServiceMock)
	server := NewServer(logger, appConfig, subscribersService, currencyService)

	request := httptest.NewRequest(http.MethodPost, "/subscribe", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err = server.SubscribeHandler(ctx)
	httpErr, ok := err.(*echo.HTTPError)

	// assert
	assert.Equal(s.T(), true, ok)
	assert.Equal(s.T(), expectedBody, httpErr.Message)
	assert.Equal(s.T(), http.StatusBadRequest, httpErr.Code)
}

func (s *RequestHandlersUnitTestSuite) TestSubscribeHandler_RecordAlreadyExists() {
	// arrange
	testSubscriber := "test_mail@gmail.com"

	expectedBody := map[string]string{
		"error": "email already exists",
	}
	appConfig, err := config.NewAppConfig(".env.test")
	assert.NoError(s.T(), err)

	logger := log.New(os.Stdout, "", appConfig.LogLevel)

	subscribersService := new(SubscribersServiceMock)
	subscribersService.On("Add", testSubscriber).Return(storage.RecordAlreadyExistsError{})

	currencyService := new(CurrencyServiceMock)
	server := NewServer(logger, appConfig, subscribersService, currencyService)

	data := url.Values{}
	data.Set("email", testSubscriber)

	request := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err = server.SubscribeHandler(ctx)
	httpErr, ok := err.(*echo.HTTPError)

	// assert
	assert.Equal(s.T(), true, ok)
	assert.Equal(s.T(), expectedBody, httpErr.Message)
	assert.Equal(s.T(), http.StatusConflict, httpErr.Code)
}

func (s *RequestHandlersUnitTestSuite) TestSubscribeHandler_InternalError() {
	// arrange
	testSubscriber := "test_mail@gmail.com"

	expectedBody := map[string]string{
		"error": "internal error",
	}
	appConfig, err := config.NewAppConfig(".env.test")
	assert.NoError(s.T(), err)

	logger := log.New(os.Stdout, "", appConfig.LogLevel)

	subscribersService := new(SubscribersServiceMock)
	subscribersService.On("Add", testSubscriber).Return(errors.New("internal error"))

	currencyService := new(CurrencyServiceMock)
	server := NewServer(logger, appConfig, subscribersService, currencyService)

	data := url.Values{}
	data.Set("email", testSubscriber)

	request := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err = server.SubscribeHandler(ctx)
	httpErr, ok := err.(*echo.HTTPError)

	// assert
	assert.Equal(s.T(), true, ok)
	assert.Equal(s.T(), expectedBody, httpErr.Message)
	assert.Equal(s.T(), http.StatusInternalServerError, httpErr.Code)
}

func (s *RequestHandlersUnitTestSuite) TestSubscribeHandler_StatusOK() {
	// arrange
	testSubscriber := "test_mail@gmail.com"

	appConfig, err := config.NewAppConfig(".env.test")
	assert.NoError(s.T(), err)

	logger := log.New(os.Stdout, "", appConfig.LogLevel)

	subscribersService := new(SubscribersServiceMock)
	subscribersService.On("Add", testSubscriber).Return(nil)

	currencyService := new(CurrencyServiceMock)
	server := NewServer(logger, appConfig, subscribersService, currencyService)

	data := url.Values{}
	data.Set("email", testSubscriber)

	request := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err = server.SubscribeHandler(ctx)

	// assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, responseRecorder.Code)
}

func (s *RequestHandlersUnitTestSuite) TestSendEmailsHandler_InternalError() {
	// arrange
	expectedBody := map[string]string{
		"error": "internal error",
	}
	rate := 50000.0
	fromCurrency := "BTC"
	toCurrency := "UAH"

	appConfig, err := config.NewAppConfig(".env.test")
	assert.NoError(s.T(), err)

	logger := log.New(os.Stdout, "", appConfig.LogLevel)

	subscribersService := new(SubscribersServiceMock)
	subscribersService.On("SendEmails", rate, fromCurrency, toCurrency).Return(nil)

	currencyService := new(CurrencyServiceMock)
	currencyService.On("GetCurrencyRate", currency.Btc, currency.Uah).Return(0.0, errors.New("internal error"))

	server := NewServer(logger, appConfig, subscribersService, currencyService)

	request := httptest.NewRequest(http.MethodPost, "/sendEmails", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err = server.SendEmailsHandler(ctx)
	httpErr, ok := err.(*echo.HTTPError)

	// assert
	assert.Equal(s.T(), true, ok)
	assert.Equal(s.T(), expectedBody, httpErr.Message)
	assert.Equal(s.T(), http.StatusInternalServerError, httpErr.Code)
}

func (s *RequestHandlersUnitTestSuite) TestSendEmailsHandler_SendMailError() {
	// arrange
	testSubscriber := "test_mail@gmail.com"

	expectedBody := map[string][]string{
		"failedEmails": {testSubscriber},
	}
	rate := 50000.0
	fromCurrency := "BTC"
	toCurrency := "UAH"

	appConfig, err := config.NewAppConfig(".env.test")
	assert.NoError(s.T(), err)

	logger := log.New(os.Stdout, "", appConfig.LogLevel)

	subscribersService := new(SubscribersServiceMock)
	subscribersService.On("SendEmails", rate, fromCurrency, toCurrency).
		Return(subscribers.SendMailError{
			Subscribers: []string{testSubscriber},
		})

	currencyService := new(CurrencyServiceMock)
	currencyService.On("GetCurrencyRate", currency.Btc, currency.Uah).Return(50000.0, nil)

	server := NewServer(logger, appConfig, subscribersService, currencyService)

	request := httptest.NewRequest(http.MethodPost, "/sendEmails", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err = server.SendEmailsHandler(ctx)

	// assert
	assert.NoError(s.T(), err)

	var actualBody map[string][]string
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &actualBody)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), expectedBody, actualBody)
	assert.Equal(s.T(), http.StatusOK, responseRecorder.Code)
}

func (s *RequestHandlersUnitTestSuite) TestSendEmailsHandler_StatusOK() {
	// arrange
	rate := 50000.0
	fromCurrency := "BTC"
	toCurrency := "UAH"

	appConfig, err := config.NewAppConfig(".env.test")
	assert.NoError(s.T(), err)

	logger := log.New(os.Stdout, "", appConfig.LogLevel)

	subscribersService := new(SubscribersServiceMock)
	subscribersService.On("SendEmails", rate, fromCurrency, toCurrency).Return(nil)

	currencyService := new(CurrencyServiceMock)
	currencyService.On("GetCurrencyRate", currency.Btc, currency.Uah).Return(50000.0, nil)

	server := NewServer(logger, appConfig, subscribersService, currencyService)

	request := httptest.NewRequest(http.MethodPost, "/sendEmails", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err = server.SendEmailsHandler(ctx)

	// assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, responseRecorder.Code)
}

func TestRequestHandlersUnitTestSuite(t *testing.T) {
	suite.Run(t, new(RequestHandlersUnitTestSuite))
}
