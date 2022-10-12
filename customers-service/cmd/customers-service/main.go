package main

import (
	"context"
	"customers-service/config"
	"customers-service/internal/api"
	"customers-service/internal/customers/controller"
	"customers-service/internal/customers/domain"
	"customers-service/internal/customers/infrastucture/postgres"
	"customers-service/pkg/postgres_client"
	"github.com/labstack/gommon/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	appConfig, err := config.NewAppConfig(".env")
	if err != nil {
		panic(err)
	}

	logger := log.New("")
	logger.SetLevel(log.Lvl(appConfig.LogLevel))

	postgresClient, err := postgres_client.NewPostgresClient(appConfig.PostgresURL)
	if err != nil {
		panic(err)
	}

	repo := postgres.NewCustomersRepo(postgresClient)
	service := domain.NewCustomersService(logger, repo)
	controller := controller.NewCustomersController(appConfig, logger, service)
	server := api.NewServer(logger, appConfig, controller)

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		logger.Infof("Server started")
		server.Run()
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	waitGroup.Add(1)
	go func() {
		<-signals
		defer waitGroup.Done()

		if err := server.Shutdown(context.Background()); err != nil {
			logger.Errorf("Server stopped with error: %s", err)
			return
		}
		logger.Infof("Server stopped gracefully")

		if err := postgres_client.ClosePostgresClient(postgresClient); err != nil {
			logger.Errorf("Failed to close postgres client: %s", err)
			return
		}
	}()

	waitGroup.Wait()
}
