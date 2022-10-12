package api

import (
	"context"
	"customers-service/config"
	"customers-service/internal/common"
	"customers-service/internal/customers/controller"
	"errors"
	"github.com/dtm-labs/dtm/dtmutil"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

type Server struct {
	logger              common.Logger
	config              *config.AppConfig
	customersController *controller.CustomersController
	httpServer          *http.Server
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
	router.POST("/customers", s.customersController.CreateCustomerHandler)
	router.DELETE("/customers/:id", s.customersController.DeleteCustomerHandler)

	// internal endpoints
	router.POST("/register-customer",
		dtmutil.WrapHandler2(s.customersController.RegisterCustomer))
	router.POST("/register-customer-compensate",
		dtmutil.WrapHandler2(s.customersController.RegisterCustomerCompensate))

	return router
}

func NewServer(
	logger common.Logger,
	cfg *config.AppConfig,
	customersController *controller.CustomersController,
) *Server {
	return &Server{
		logger:              logger,
		config:              cfg,
		customersController: customersController,
		httpServer:          &http.Server{},
	}
}
