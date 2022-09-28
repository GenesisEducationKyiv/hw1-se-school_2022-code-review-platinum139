package api

import (
	"bitcoin-service/config"
	currency "bitcoin-service/internal/currency/controller"
	notification "bitcoin-service/internal/notification/controller"
	subscribers "bitcoin-service/internal/subscribers/controller"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Server struct {
	logger                 *log.Logger
	config                 *config.AppConfig
	notificationController *notification.Controller
	currencyController     *currency.Controller
	subscribersController  *subscribers.Controller
	echoServer             *echo.Echo
}

func (s Server) RegisterRoutes() {
	s.echoServer.GET("/rate", s.currencyController.GetRateHandler)
	s.echoServer.POST("/subscribe", s.subscribersController.SubscribeHandler)
	s.echoServer.POST("/sendEmails", s.notificationController.SendEmailsHandler)
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
	logger *log.Logger,
	cfg *config.AppConfig,
	notificationController *notification.Controller,
	currencyController *currency.Controller,
	subscribersController *subscribers.Controller,
) *Server {
	echoServer := echo.New()
	return &Server{
		logger:                 logger,
		config:                 cfg,
		notificationController: notificationController,
		currencyController:     currencyController,
		subscribersController:  subscribersController,
		echoServer:             echoServer,
	}
}
