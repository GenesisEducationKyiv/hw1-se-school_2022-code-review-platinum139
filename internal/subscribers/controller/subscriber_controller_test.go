package controller

import (
	"bitcoin-service/config"
	"bitcoin-service/internal/subscribers/domain"
	"bitcoin-service/pkg/file_storage"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SubscribersServiceMock struct {
	mock.Mock
}

func (s *SubscribersServiceMock) Subscribe(subscriber domain.Subscriber) error {
	args := s.Called(subscriber)
	return args.Error(0)
}

func (s *SubscribersServiceMock) GetSubscribers() ([]domain.Subscriber, error) {
	args := s.Called()
	return args.Get(0).([]domain.Subscriber), args.Error(1)
}

type SubscribersControllerTestSuite struct {
	suite.Suite
	config *config.AppConfig
	logger *log.Logger
}

func (s *SubscribersControllerTestSuite) SetupSuite() {
	appConfig, err := config.NewAppConfig(".env.test")
	assert.NoError(s.T(), err)
	s.config = appConfig

	s.logger = log.New(os.Stdout, "", appConfig.LogLevel)
}

func (s *SubscribersControllerTestSuite) TestSubscribeHandler_EmailNotProvided() {
	// arrange
	expectedBody := map[string]string{
		"error": "email must be provided",
	}

	subscribersService := new(SubscribersServiceMock)
	controller := NewSubscribersController(s.logger, s.config, subscribersService)

	request := httptest.NewRequest(http.MethodPost, "/subscribe", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := controller.SubscribeHandler(ctx)
	httpErr, ok := err.(*echo.HTTPError)

	// assert
	assert.Equal(s.T(), true, ok)
	assert.Equal(s.T(), expectedBody, httpErr.Message)
	assert.Equal(s.T(), http.StatusBadRequest, httpErr.Code)
}

func (s *SubscribersControllerTestSuite) TestSubscribeHandler_RecordAlreadyExists() {
	// arrange
	testSubscriber := domain.Subscriber{
		Email: "test_mail@gmail.com",
	}

	expectedBody := map[string]string{
		"error": "email already exists",
	}

	subscribersService := new(SubscribersServiceMock)
	subscribersService.On("Subscribe", testSubscriber).Return(file_storage.RecordAlreadyExistsError{})

	controller := NewSubscribersController(s.logger, s.config, subscribersService)

	data := url.Values{}
	data.Set("email", testSubscriber.Email)

	request := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := controller.SubscribeHandler(ctx)
	httpErr, ok := err.(*echo.HTTPError)

	// assert
	assert.Equal(s.T(), true, ok)
	assert.Equal(s.T(), expectedBody, httpErr.Message)
	assert.Equal(s.T(), http.StatusConflict, httpErr.Code)
}

func (s *SubscribersControllerTestSuite) TestSubscribeHandler_InternalError() {
	// arrange
	testSubscriber := domain.Subscriber{
		Email: "test_mail@gmail.com",
	}

	expectedBody := map[string]string{
		"error": "internal error",
	}

	subscribersService := new(SubscribersServiceMock)
	subscribersService.On("Subscribe", testSubscriber).Return(errors.New("internal error"))

	controller := NewSubscribersController(s.logger, s.config, subscribersService)

	data := url.Values{}
	data.Set("email", testSubscriber.Email)

	request := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := controller.SubscribeHandler(ctx)
	httpErr, ok := err.(*echo.HTTPError)

	// assert
	assert.Equal(s.T(), true, ok)
	assert.Equal(s.T(), expectedBody, httpErr.Message)
	assert.Equal(s.T(), http.StatusInternalServerError, httpErr.Code)
}

func (s *SubscribersControllerTestSuite) TestSubscribeHandler_StatusOK() {
	// arrange
	testSubscriber := domain.Subscriber{
		Email: "test_mail@gmail.com",
	}

	subscribersService := new(SubscribersServiceMock)
	subscribersService.On("Subscribe", testSubscriber).Return(nil)

	controller := NewSubscribersController(s.logger, s.config, subscribersService)

	data := url.Values{}
	data.Set("email", testSubscriber.Email)

	request := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := controller.SubscribeHandler(ctx)

	// assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, responseRecorder.Code)
}

func TestSubscribersControllerTestSuite(t *testing.T) {
	suite.Run(t, new(SubscribersControllerTestSuite))
}
