package api

import (
	"bitcoin-service/config"
	"bitcoin-service/internal/subscribers"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type Server struct {
	logger     *log.Logger
	config     *config.AppConfig
	service    *subscribers.Service
	echoServer *echo.Echo
}

func (s Server) RegisterRoutes() {
	s.echoServer.GET("/rate", s.GetRate)
	s.echoServer.POST("/subscribe", s.Subscribe)
	s.echoServer.POST("/sendEmails", s.SendEmails)
}

func (s Server) Run() {
	addr := fmt.Sprintf("%s:%s", s.config.ServerHost, s.config.ServerPort)
	if err := s.echoServer.Start(addr); err != nil && err != http.ErrServerClosed {
		s.logger.Print("Server stopped with error:", err)
	}
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.echoServer.Shutdown(ctx)
}

func NewServer(logger *log.Logger, cfg *config.AppConfig, service *subscribers.Service) *Server {
	echoServer := echo.New()

	return &Server{
		logger:     logger,
		config:     cfg,
		service:    service,
		echoServer: echoServer,
	}
}
