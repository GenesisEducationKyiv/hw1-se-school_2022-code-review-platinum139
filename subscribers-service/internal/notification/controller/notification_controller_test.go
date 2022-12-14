package controller

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"subscribers-service/config"
	"subscribers-service/internal/common"
	notificationDomain "subscribers-service/internal/notification/domain"
	subscribersDomain "subscribers-service/internal/subscribers/domain"
	"testing"
)

type NotificationServiceMock struct {
	mock.Mock
}

func (s *NotificationServiceMock) Notify() error {
	args := s.Called()
	return args.Error(0)
}

type NotificationServiceTestSuite struct {
	suite.Suite
	config *config.AppConfig
	logger common.Logger
}

func (s *NotificationServiceTestSuite) SetupSuite() {
	appConfig, err := config.NewAppConfig(".env")
	assert.NoError(s.T(), err)
	s.config = appConfig

	s.logger = log.New("")
}

func (s *NotificationServiceTestSuite) TestSendEmailsHandler_InternalError() {
	// arrange
	expectedBody := map[string]string{
		"error": "internal error",
	}

	notificationService := new(NotificationServiceMock)
	notificationService.On("Notify").Return(errors.New("rate not found"))

	controller := NewNotificationController(s.logger, notificationService)

	request := httptest.NewRequest(http.MethodPost, "/sendEmails", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := controller.SendEmailsHandler(ctx)
	httpErr, ok := err.(*echo.HTTPError)

	// assert
	assert.Equal(s.T(), true, ok)
	assert.Equal(s.T(), expectedBody, httpErr.Message)
	assert.Equal(s.T(), http.StatusInternalServerError, httpErr.Code)
}

func (s *NotificationServiceTestSuite) TestSendEmailsHandler_SendMailError() {
	// arrange
	testSubscriber := subscribersDomain.Subscriber{
		Email: "test_mail@gmail.com",
	}

	expectedBody := map[string][]string{
		"failedEmails": {testSubscriber.Email},
	}

	notificationService := new(NotificationServiceMock)
	notificationService.On("Notify").
		Return(notificationDomain.SendMessageError{
			FailedSubscribers: []string{testSubscriber.Email},
		})

	controller := NewNotificationController(s.logger, notificationService)

	request := httptest.NewRequest(http.MethodPost, "/sendEmails", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := controller.SendEmailsHandler(ctx)

	// assert
	assert.NoError(s.T(), err)

	var actualBody map[string][]string
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &actualBody)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), expectedBody, actualBody)
	assert.Equal(s.T(), http.StatusOK, responseRecorder.Code)
}

func (s *NotificationServiceTestSuite) TestSendEmailsHandler_StatusOK() {
	// arrange
	notificationService := new(NotificationServiceMock)
	notificationService.On("Notify").Return(nil)

	controller := NewNotificationController(s.logger, notificationService)

	request := httptest.NewRequest(http.MethodPost, "/sendEmails", nil)
	responseRecorder := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(request, responseRecorder)

	// act
	err := controller.SendEmailsHandler(ctx)

	// assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, responseRecorder.Code)
}

func TestNotificationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationServiceTestSuite))
}
