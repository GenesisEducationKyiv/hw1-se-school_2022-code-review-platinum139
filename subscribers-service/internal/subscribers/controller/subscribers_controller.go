package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"subscribers-service/config"
	"subscribers-service/internal/common"
	"subscribers-service/internal/subscribers/domain"
	"subscribers-service/pkg/file_storage"
)

type SubscribeRequest struct {
	Email string `form:"email"`
}

type Controller struct {
	logger             common.Logger
	config             *config.AppConfig
	subscribersService domain.SubscribersService
}

func (c Controller) SubscribeHandler(context echo.Context) error {
	c.logger.Debugf("SubscribeHandler started")

	subscribeRequest := new(SubscribeRequest)
	if err := context.Bind(subscribeRequest); err != nil {
		c.logger.Errorf("Unable to get form data from request:", err)
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "bad request",
		})
	}

	if subscribeRequest.Email == "" {
		c.logger.Errorf("Bad request: email was not provided")
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "email must be provided",
		})
	}

	subscriber := domain.Subscriber{
		Email: subscribeRequest.Email,
	}

	err := c.subscribersService.Subscribe(subscriber)
	if err != nil {
		if errors.Is(err, file_storage.RecordAlreadyExistsError{}) {
			c.logger.Errorf("Email already exists in storage:", subscriber.Email)
			return echo.NewHTTPError(http.StatusConflict, map[string]string{
				"error": "email already exists",
			})
		}
		c.logger.Errorf("Unable to add subscriber to storage:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"error": "internal error",
		})
	}

	return context.JSON(http.StatusOK, "")
}

func NewSubscribersController(
	logger common.Logger,
	config *config.AppConfig,
	subscribersService domain.SubscribersService,
) *Controller {
	return &Controller{
		logger:             logger,
		config:             config,
		subscribersService: subscribersService,
	}
}
