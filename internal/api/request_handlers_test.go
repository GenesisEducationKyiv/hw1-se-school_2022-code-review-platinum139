package api

import (
	"bitcoin-service/config"
	"bitcoin-service/internal/common"
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

func (s *SubscribersServiceMock) Subscribe(subscriber subscribers.Subscriber) error {
	args := s.Called(subscriber)
	return args.Error(0)
}

func (s *SubscribersServiceMock) SendEmails(message common.Message) error {
	args := s.Called(message)
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
	config *config.AppConfig
	logger *log.Logger
}

func (s *RequestHandlersUnitTestSuite) SetupSuite() {
	appConfig, err := config.NewAppConfig(".env.test")
	assert.NoError(s.T(), err)
	s.config = appConfig

	s.logger = log.New(os.Stdout, "", appConfig.LogLevel)
}

func (s *RequestHandlersUnitTestSuite) TestGetRateHandler_StatusOK() {
	// arrange
	subscribersService := new(SubscribersServiceMock)

	currencyService := new(CurrencyServiceMock)
	currencyService.On("GetCurrencyRate", currency.Btc, currency.Uah).Return(50000.0, nil)

	server := NewServer(s.logger, s.config, subscribersService, currencyService)

	request := httptest.NewRequest(http.MethodGet, "/rate", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := server.GetRateHandler(ctx)

	// assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, responseRecorder.Code)
}

func (s *RequestHandlersUnitTestSuite) TestGetRateHandler_StatusBadRequest() {
	// arrange
	expectedBody := map[string]string{
		"error": "invalid status value",
	}

	subscribersService := new(SubscribersServiceMock)

	currencyService := new(CurrencyServiceMock)
	currencyService.On("GetCurrencyRate", currency.Btc, currency.Uah).Return(0.0, http.ErrServerClosed)

	server := NewServer(s.logger, s.config, subscribersService, currencyService)

	request := httptest.NewRequest(http.MethodGet, "/rate", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := server.GetRateHandler(ctx)
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

	subscribersService := new(SubscribersServiceMock)
	currencyService := new(CurrencyServiceMock)
	server := NewServer(s.logger, s.config, subscribersService, currencyService)

	request := httptest.NewRequest(http.MethodPost, "/subscribe", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := server.SubscribeHandler(ctx)
	httpErr, ok := err.(*echo.HTTPError)

	// assert
	assert.Equal(s.T(), true, ok)
	assert.Equal(s.T(), expectedBody, httpErr.Message)
	assert.Equal(s.T(), http.StatusBadRequest, httpErr.Code)
}

func (s *RequestHandlersUnitTestSuite) TestSubscribeHandler_RecordAlreadyExists() {
	// arrange
	testSubscriber := subscribers.Subscriber{
		Email: "test_mail@gmail.com",
	}

	expectedBody := map[string]string{
		"error": "email already exists",
	}

	subscribersService := new(SubscribersServiceMock)
	subscribersService.On("Subscribe", testSubscriber).Return(storage.RecordAlreadyExistsError{})

	currencyService := new(CurrencyServiceMock)
	server := NewServer(s.logger, s.config, subscribersService, currencyService)

	data := url.Values{}
	data.Set("email", testSubscriber.Email)

	request := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := server.SubscribeHandler(ctx)
	httpErr, ok := err.(*echo.HTTPError)

	// assert
	assert.Equal(s.T(), true, ok)
	assert.Equal(s.T(), expectedBody, httpErr.Message)
	assert.Equal(s.T(), http.StatusConflict, httpErr.Code)
}

func (s *RequestHandlersUnitTestSuite) TestSubscribeHandler_InternalError() {
	// arrange
	testSubscriber := subscribers.Subscriber{
		Email: "test_mail@gmail.com",
	}

	expectedBody := map[string]string{
		"error": "internal error",
	}

	subscribersService := new(SubscribersServiceMock)
	subscribersService.On("Subscribe", testSubscriber).Return(errors.New("internal error"))

	currencyService := new(CurrencyServiceMock)
	server := NewServer(s.logger, s.config, subscribersService, currencyService)

	data := url.Values{}
	data.Set("email", testSubscriber.Email)

	request := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := server.SubscribeHandler(ctx)
	httpErr, ok := err.(*echo.HTTPError)

	// assert
	assert.Equal(s.T(), true, ok)
	assert.Equal(s.T(), expectedBody, httpErr.Message)
	assert.Equal(s.T(), http.StatusInternalServerError, httpErr.Code)
}

func (s *RequestHandlersUnitTestSuite) TestSubscribeHandler_StatusOK() {
	// arrange
	testSubscriber := subscribers.Subscriber{
		Email: "test_mail@gmail.com",
	}

	subscribersService := new(SubscribersServiceMock)
	subscribersService.On("Subscribe", testSubscriber).Return(nil)

	currencyService := new(CurrencyServiceMock)
	server := NewServer(s.logger, s.config, subscribersService, currencyService)

	data := url.Values{}
	data.Set("email", testSubscriber.Email)

	request := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := server.SubscribeHandler(ctx)

	// assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, responseRecorder.Code)
}

func (s *RequestHandlersUnitTestSuite) TestSendEmailsHandler_InternalError() {
	// arrange
	expectedBody := map[string]string{
		"error": "internal error",
	}

	message := common.Message{
		Subject: "BTC to UAH rate",
		Text:    "Rate = 50000.00",
	}

	subscribersService := new(SubscribersServiceMock)
	subscribersService.On("SendEmails", message).Return(nil)

	currencyService := new(CurrencyServiceMock)
	currencyService.On("GetCurrencyRate", currency.Btc, currency.Uah).Return(0.0, errors.New("internal error"))

	server := NewServer(s.logger, s.config, subscribersService, currencyService)

	request := httptest.NewRequest(http.MethodPost, "/sendEmails", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := server.SendEmailsHandler(ctx)
	httpErr, ok := err.(*echo.HTTPError)

	// assert
	assert.Equal(s.T(), true, ok)
	assert.Equal(s.T(), expectedBody, httpErr.Message)
	assert.Equal(s.T(), http.StatusInternalServerError, httpErr.Code)
}

func (s *RequestHandlersUnitTestSuite) TestSendEmailsHandler_SendMailError() {
	// arrange
	rate := 50000.0

	testSubscriber := subscribers.Subscriber{
		Email: "test_mail@gmail.com",
	}

	expectedBody := map[string][]string{
		"failedEmails": {testSubscriber.Email},
	}

	message := common.Message{
		Subject: "BTC to UAH rate",
		Text:    "Rate = 50000.00",
	}

	subscribersService := new(SubscribersServiceMock)
	subscribersService.On("SendEmails", message).
		Return(subscribers.SendMessageError{
			FailedSubscribers: []string{testSubscriber.Email},
		})

	currencyService := new(CurrencyServiceMock)
	currencyService.On("GetCurrencyRate", currency.Btc, currency.Uah).Return(rate, nil)

	server := NewServer(s.logger, s.config, subscribersService, currencyService)

	request := httptest.NewRequest(http.MethodPost, "/sendEmails", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := server.SendEmailsHandler(ctx)

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

	message := common.Message{
		Subject: "BTC to UAH rate",
		Text:    "Rate = 50000.00",
	}

	subscribersService := new(SubscribersServiceMock)
	subscribersService.On("SendEmails", message).Return(nil)

	currencyService := new(CurrencyServiceMock)
	currencyService.On("GetCurrencyRate", currency.Btc, currency.Uah).Return(rate, nil)

	server := NewServer(s.logger, s.config, subscribersService, currencyService)

	request := httptest.NewRequest(http.MethodPost, "/sendEmails", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := server.SendEmailsHandler(ctx)

	// assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, responseRecorder.Code)
}

func TestRequestHandlersUnitTestSuite(t *testing.T) {
	suite.Run(t, new(RequestHandlersUnitTestSuite))
}
