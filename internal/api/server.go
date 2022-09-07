package api

import (
	"bitcoin-service/config"
	"bitcoin-service/pkg/currency"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SubscribersService interface {
	Add(subscriber string) error
	SendEmails(rate float64, fromCurrency, toCurrency string) error
}

type CurrencyService interface {
	GetCurrencyRate(from currency.Currency, to currency.Currency) (float64, error)
}

type Server struct {
	logger             *log.Logger
	config             *config.AppConfig
	subscribersService SubscribersService
	currencyService    CurrencyService
	echoServer         *echo.Echo
}

func (s Server) RegisterRoutes() {
	s.echoServer.GET("/rate", s.GetRateHandler)
	s.echoServer.POST("/subscribe", s.SubscribeHandler)
	s.echoServer.POST("/sendEmails", s.SendEmailsHandler)
}

func (s Server) Run() {
	addr := fmt.Sprintf("%s:%s", s.config.ServerHost, s.config.ServerPort)
	if err := s.echoServer.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Print("Server stopped with error:", err)
	}
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.echoServer.Shutdown(ctx)
}

func NewServer(
	logger *log.Logger, cfg *config.AppConfig,
	subscribersService SubscribersService,
	currencyService CurrencyService) *Server {

	echoServer := echo.New()

	return &Server{
		logger:             logger,
		config:             cfg,
		subscribersService: subscribersService,
		currencyService:    currencyService,
		echoServer:         echoServer,
	}
}
