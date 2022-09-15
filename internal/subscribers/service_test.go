package subscribers

import (
	"bitcoin-service/config"
	"bitcoin-service/internal/common"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SubscribersRepoMock struct {
	mock.Mock
}

func (s *SubscribersRepoMock) CreateSubscriber(subscriber Subscriber) error {
	args := s.Called(subscriber)
	return args.Error(0)
}

func (s *SubscribersRepoMock) GetSubscribers() ([]Subscriber, error) {
	args := s.Called()
	return args.Get(0).([]Subscriber), args.Error(1)
}

type EmailSenderMock struct {
	mock.Mock
}

func (s *EmailSenderMock) Send(receiverEmail string, subject string, text string) error {
	args := s.Called(receiverEmail, subject, text)
	return args.Error(0)
}

type SubscribersServiceUnitTestSuite struct {
	suite.Suite
	logger *log.Logger
}

func (s *SubscribersServiceUnitTestSuite) SetupSuite() {
	appConfig, err := config.NewAppConfig(".env.test")
	assert.NoError(s.T(), err)

	s.logger = log.New(os.Stdout, "", appConfig.LogLevel)
}

func (s *SubscribersServiceUnitTestSuite) TestSendEmails_Positive() {
	// arrange
	receiver := "test_mail@gmail.com"
	subscribers := []Subscriber{{Email: receiver}}

	message := common.Message{
		Subject: "BTC to UAH rate",
		Text:    "Rate = 50000.00",
	}

	subscribersRepo := new(SubscribersRepoMock)
	subscribersRepo.On("GetSubscribers").Return(subscribers, nil)

	emailSender := new(EmailSenderMock)
	emailSender.On("Send", receiver, message.Subject, message.Text).Return(nil)

	service := NewSubscribersService(s.logger, subscribersRepo, emailSender)

	// act
	err := service.SendEmails(message)

	// assert
	assert.NoError(s.T(), err)

	subscribersRepo.AssertCalled(s.T(), "GetSubscribers")
	emailSender.AssertCalled(s.T(), "Send", receiver, message.Subject, message.Text)
}

func (s *SubscribersServiceUnitTestSuite) TestSendEmails_FailedEmails() {
	// arrange
	receiver := "test_mail@gmail.com"
	subscribers := []Subscriber{{Email: receiver}}
	failedEmails := []string{receiver}

	message := common.Message{
		Subject: "BTC to UAH rate",
		Text:    "Rate = 50000.00",
	}

	subscribersRepo := new(SubscribersRepoMock)
	subscribersRepo.On("GetSubscribers").Return(subscribers, nil)

	emailSender := new(EmailSenderMock)
	emailSender.On("Send", receiver, message.Subject, message.Text).
		Return(errors.New("failed to send"))

	service := NewSubscribersService(s.logger, subscribersRepo, emailSender)

	// act
	err := service.SendEmails(message)

	// assert
	assert.EqualError(s.T(), err, SendMessageError{FailedSubscribers: failedEmails}.Error())

	subscribersRepo.AssertCalled(s.T(), "GetSubscribers")
	emailSender.AssertCalled(s.T(), "Send", receiver, message.Subject, message.Text)
}

func TestSubscribersServiceUnitTestSuite(t *testing.T) {
	suite.Run(t, new(SubscribersServiceUnitTestSuite))
}
