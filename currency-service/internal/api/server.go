package api

import (
	"context"
	"currency-service/config"
	"currency-service/internal/currency/controller"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Server struct {
	logger             *log.Logger
	config             *config.AppConfig
	currencyController *controller.Controller
	echoServer         *echo.Echo
}

func (s Server) RegisterRoutes() {
	s.echoServer.GET("/rate", s.currencyController.GetRateHandler)
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
	currencyController *controller.Controller,
) *Server {
	echoServer := echo.New()
	return &Server{
		logger:             logger,
		config:             cfg,
		currencyController: currencyController,
		echoServer:         echoServer,
	}
}
