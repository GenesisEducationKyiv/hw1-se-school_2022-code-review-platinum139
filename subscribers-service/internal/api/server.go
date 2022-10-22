package api

import (
	"context"
	"errors"
	"github.com/dtm-labs/dtm/dtmutil"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"subscribers-service/config"
	"subscribers-service/internal/common"
	notification "subscribers-service/internal/notification/controller"
	subscribers "subscribers-service/internal/subscribers/controller"
	"time"
)

type Server struct {
	logger                 common.Logger
	config                 *config.AppConfig
	notificationController *notification.Controller
	subscribersController  *subscribers.Controller
	httpServer             *http.Server
}

func (s Server) Run() {
	addr := net.JoinHostPort(s.config.ServerHost, s.config.ServerPort)
	router := s.getRouter()

	s.httpServer.Addr = addr
	s.httpServer.Handler = router

	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Errorf("Server stopped with error:", err)
	}
	s.logger.Infof("Server started on port %s", s.config.ServerPort)
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s Server) getRouter() *gin.Engine {
	router := gin.Default()

	// public endpoints
	router.POST("/subscribe", s.subscribersController.SubscribeHandler)
	router.POST("/sendEmails", s.notificationController.SendEmailsHandler)

	// internal endpoints
	router.POST("/register-subscriber",
		dtmutil.WrapHandler2(s.subscribersController.RegisterSubscriber))
	router.POST("/register-subscriber-compensate",
		dtmutil.WrapHandler2(s.subscribersController.RegisterSubscriberCompensate))

	return router
}

func NewServer(
	logger common.Logger,
	cfg *config.AppConfig,
	notificationController *notification.Controller,
	subscribersController *subscribers.Controller,
) *Server {
	httpServer := http.Server{ReadHeaderTimeout: time.Second}
	return &Server{
		logger:                 logger,
		config:                 cfg,
		notificationController: notificationController,
		subscribersController:  subscribersController,
		httpServer:             &httpServer,
	}
}
