package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"subscribers-service/config"
	notification "subscribers-service/internal/notification/controller"
	subscribers "subscribers-service/internal/subscribers/controller"
)

type Server struct {
	logger                 *log.Logger
	config                 *config.AppConfig
	notificationController *notification.Controller
	subscribersController  *subscribers.Controller
	echoServer             *echo.Echo
}

func (s Server) RegisterRoutes() {
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
	subscribersController *subscribers.Controller,
) *Server {
	echoServer := echo.New()
	return &Server{
		logger:                 logger,
		config:                 cfg,
		notificationController: notificationController,
		subscribersController:  subscribersController,
		echoServer:             echoServer,
	}
}
