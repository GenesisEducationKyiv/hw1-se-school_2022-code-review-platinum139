package controller

import (
	"currency-service/internal/currency/domain"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
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
