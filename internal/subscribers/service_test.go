package subscribers

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type StorageMock struct {
	mock.Mock
}

func (s *StorageMock) Add(record string) error {
	args := s.Called(record)
	return args.Error(0)
}

func (s *StorageMock) GetAll() ([]string, error) {
	args := s.Called()
	return args.Get(0).([]string), args.Error(1)
}

type EmailSenderMock struct {
	mock.Mock
}

func (s *EmailSenderMock) SendEmail(receiverEmail string, subject string, text string) error {
	args := s.Called(receiverEmail, subject, text)
	return args.Error(0)
}

type SubscribersServiceUnitTestSuite struct {
	suite.Suite
}

func (s *SubscribersServiceUnitTestSuite) TestSendEmails_Positive() {
	// arrange
	rate := 50000.0
	fromCurrency := "BTC"
	toCurrency := "UAH"
	receiver := "test_mail@gmail.com"
	subscribers := []string{receiver}
	subject := "BTC to UAH rate"
	text := "Rate = 50000.00"

	storage := new(StorageMock)
	storage.On("GetAll").Return(subscribers, nil)

	emailSender := new(EmailSenderMock)
	emailSender.On("SendEmail", receiver, subject, text).Return(nil)

	logger := log.New(os.Stdout, "", 4)
	service := NewSubscribersService(logger, storage, emailSender)

	// act
	err := service.SendEmails(rate, fromCurrency, toCurrency)

	// assert
	assert.NoError(s.T(), err)

	storage.AssertCalled(s.T(), "GetAll")
	emailSender.AssertCalled(s.T(), "SendEmail", receiver, subject, text)
}

func (s *SubscribersServiceUnitTestSuite) TestSendEmails_FailedEmails() {
	// arrange
	rate := 50000.0
	fromCurrency := "BTC"
	toCurrency := "UAH"
	receiver := "test_mail@gmail.com"
	subscribers := []string{receiver}
	subject := "BTC to UAH rate"
	text := "Rate = 50000.00"

	storage := new(StorageMock)
	storage.On("GetAll").Return(subscribers, nil)

	emailSender := new(EmailSenderMock)
	emailSender.On("SendEmail", receiver, subject, text).Return(errors.New("failed to send"))

	logger := log.New(os.Stdout, "", 4)
	service := NewSubscribersService(logger, storage, emailSender)

	// act
	err := service.SendEmails(rate, fromCurrency, toCurrency)

	// assert
	assert.EqualError(s.T(), err, SendMailError{Subscribers: subscribers}.Error())

	storage.AssertCalled(s.T(), "GetAll")
	emailSender.AssertCalled(s.T(), "SendEmail", receiver, subject, text)
}

func TestSubscribersServiceUnitTestSuite(t *testing.T) {
	suite.Run(t, new(SubscribersServiceUnitTestSuite))
}
