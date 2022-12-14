package emails

import (
	"subscribers-service/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EmailsServiceIntegrationTestSuite struct {
	suite.Suite
}

func (s *EmailsServiceIntegrationTestSuite) TestSendEmail_Positive() {
	// arrange
	appConfig, err := config.NewAppConfig(".env")
	assert.NoError(s.T(), err)

	receiverEmail := appConfig.SMTPUsername
	subject := "Test message"
	text := "This is a test message."

	emailService := NewEmailService(appConfig)

	// act
	err = emailService.Send(receiverEmail, subject, text)

	// assert
	assert.NoError(s.T(), err)
}

func TestEmailsIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(EmailsServiceIntegrationTestSuite))
}
