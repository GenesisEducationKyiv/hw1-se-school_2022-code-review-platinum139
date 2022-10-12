package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"subscribers-service/internal/common"
	"subscribers-service/internal/notification/domain"
)

type Controller struct {
	logger              common.Logger
	notificationService domain.NotificationService
}

func (c Controller) SendEmailsHandler(context *gin.Context) {
	c.logger.Debugf("SendEmailHandler started")

	err := c.notificationService.Notify()
	if err != nil {
		var sendMailErr domain.SendMessageError
		if errors.As(err, &sendMailErr) {
			c.logger.Errorf("Unable to send emails: ", sendMailErr)
			context.JSON(http.StatusOK, map[string][]string{
				"failedEmails": sendMailErr.FailedSubscribers,
			})
			return
		}
		context.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error",
		})
		return
	}

	context.Status(http.StatusOK)
}

func NewNotificationController(
	logger common.Logger,
	notificationService domain.NotificationService,
) *Controller {
	return &Controller{
		logger:              logger,
		notificationService: notificationService,
	}
}
