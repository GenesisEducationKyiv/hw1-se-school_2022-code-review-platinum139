package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"subscribers-service/internal/notification/domain"
)

type Controller struct {
	logger              *log.Logger
	notificationService domain.NotificationService
}

func (c Controller) SendEmailsHandler(context echo.Context) error {
	err := c.notificationService.Notify()
	if err != nil {
		var sendMailErr domain.SendMessageError
		if errors.As(err, &sendMailErr) {
			c.logger.Print("Unable to send emails: ", sendMailErr)
			return context.JSON(http.StatusOK, map[string][]string{
				"failedEmails": sendMailErr.FailedSubscribers,
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"error": "internal error",
		})
	}

	return context.JSON(http.StatusOK, "")
}

func NewNotificationController(
	logger *log.Logger,
	notificationService domain.NotificationService,
) *Controller {
	return &Controller{
		logger:              logger,
		notificationService: notificationService,
	}
}
