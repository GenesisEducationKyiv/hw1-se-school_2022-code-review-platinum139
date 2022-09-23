package controller

import (
	"bitcoin-service/internal/currency/domain"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger          *log.Logger
	currencyService domain.CurrencyService
}

func (c *Controller) GetRateHandler(context echo.Context) error {
	rate, err := c.currencyService.GetCurrencyRate(domain.Btc, domain.Uah)
	if err != nil {
		c.logger.Print("Unable to get currency rate:", err)
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "invalid status value",
		})
	}
	return context.JSON(http.StatusOK, rate)
}

func NewCurrencyController(logger *log.Logger, service *domain.CurrencyServiceImpl) *Controller {
	return &Controller{
		logger:          logger,
		currencyService: service,
	}
}
