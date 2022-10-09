package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"subscribers-service/config"
	"subscribers-service/internal/common"
	notification "subscribers-service/internal/notification/controller"
	subscribers "subscribers-service/internal/subscribers/controller"
)

type Server struct {
	logger                 common.Logger
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
		s.logger.Errorf("Server stopped with error:", err)
	}
	s.logger.Infof("Server started on port %s", s.config.ServerPort)
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.echoServer.Shutdown(ctx)
}

func NewServer(
	logger common.Logger,
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
