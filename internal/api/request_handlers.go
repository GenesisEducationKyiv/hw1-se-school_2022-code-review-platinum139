package api

import (
	"bitcoin-service/internal/subscribers"
	"bitcoin-service/pkg/rates"
	"bitcoin-service/pkg/storage"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	uah = "UAH"
	btc = "BTC"
)

type Subscriber struct {
	Email string `form:"email"`
}

func (s Server) GetRate(context echo.Context) error {
	rate, err := rates.GetCurrencyRate(btc, uah)
	if err != nil {
		s.logger.Print("Unable to get currency rate:", err)
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "invalid status value",
		})
	}
	return context.JSON(http.StatusOK, rate)
}

func (s Server) Subscribe(context echo.Context) error {
	subscriber := new(Subscriber)
	if err := context.Bind(subscriber); err != nil {
		s.logger.Print("Unable to get form data from request:", err)
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "bad request",
		})
	}
	if subscriber.Email == "" {
		s.logger.Print("Bad request: email was not provided")
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "email must be provided",
		})
	}

	err := s.service.Add(subscriber.Email)
	if err != nil {
		if errors.Is(err, storage.RecordAlreadyExistsError{}) {
			s.logger.Print("Email already exists in storage:", subscriber.Email)
			return echo.NewHTTPError(http.StatusConflict, map[string]string{
				"error": "email already exists",
			})
		}
		s.logger.Print("Unable to add subscriber to storage:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"error": "internal error",
		})
	}

	return context.JSON(http.StatusOK, "")
}

func (s Server) SendEmails(context echo.Context) error {
	if err := s.service.SendEmails(); err != nil {
		var sendMailErr subscribers.SendMailError
		if errors.As(err, &sendMailErr) {
			s.logger.Print("Unable to send emails: ", sendMailErr)
			return context.JSON(http.StatusOK, map[string][]string{
				"failedEmails": sendMailErr.Subscribers,
			})
		}
	}
	return context.JSON(http.StatusOK, "")
}
