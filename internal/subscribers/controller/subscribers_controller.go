package controller

import (
	"bitcoin-service/config"
	"bitcoin-service/internal/subscribers/domain"
	"bitcoin-service/pkg/file_storage"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SubscribeRequest struct {
	Email string `form:"email"`
}

type Controller struct {
	logger             *log.Logger
	config             *config.AppConfig
	subscribersService domain.SubscribersService
}

func (c Controller) SubscribeHandler(context echo.Context) error {
	subscribeRequest := new(SubscribeRequest)
	if err := context.Bind(subscribeRequest); err != nil {
		c.logger.Print("Unable to get form data from request:", err)
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "bad request",
		})
	}

	if subscribeRequest.Email == "" {
		c.logger.Print("Bad request: email was not provided")
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
			c.logger.Print("Email already exists in storage:", subscriber.Email)
			return echo.NewHTTPError(http.StatusConflict, map[string]string{
				"error": "email already exists",
			})
		}
		c.logger.Print("Unable to add subscriber to storage:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"error": "internal error",
		})
	}

	return context.JSON(http.StatusOK, "")
}

func NewSubscribersController(
	logger *log.Logger,
	config *config.AppConfig,
	subscribersService domain.SubscribersService,
) *Controller {
	return &Controller{
		logger:             logger,
		config:             config,
		subscribersService: subscribersService,
	}
}
