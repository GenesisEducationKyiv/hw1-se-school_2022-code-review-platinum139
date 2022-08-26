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

func (s Server) GetRate(c echo.Context) error {
	rate, err := rates.GetCurrencyRate(btc, uah)
	if err != nil {
		s.logger.Print("Unable to get currency rate:", err)
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "invalid status value",
		})
	}
	return c.JSON(http.StatusOK, rate)
}

func (s Server) Subscribe(c echo.Context) error {
	subscriber := new(Subscriber)
	if err := c.Bind(subscriber); err != nil {
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

	return c.JSON(http.StatusOK, "")
}

func (s Server) SendEmails(c echo.Context) error {
	if err := s.service.SendEmails(); err != nil {

		var sendMailErr subscribers.SendMailError
		if errors.As(err, &sendMailErr) {
			s.logger.Print("Unable to send emails: ", sendMailErr)
			return c.JSON(http.StatusOK, map[string][]string{
				"failedEmails": sendMailErr.Subscribers,
			})
		}
	}
	return c.JSON(http.StatusOK, "")
}
