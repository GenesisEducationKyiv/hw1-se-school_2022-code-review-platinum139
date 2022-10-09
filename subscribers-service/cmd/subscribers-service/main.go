package main

import (
	"context"
	"os"
	"os/signal"
	"subscribers-service/config"
	"subscribers-service/internal/common"
	"subscribers-service/pkg/message_brokers"
	"sync"
	"syscall"
)

const logsExchangeName = "logs"

func main() {
	appConfig, err := config.NewAppConfig(".env")
	if err != nil {
		panic(err)
	}

	rabbitMqClient, err := message_brokers.NewRabbitMqClient(message_brokers.RabbitMqConfig{
		Host:     appConfig.RabbitMqHost,
		Port:     appConfig.RabbitMqPort,
		User:     appConfig.RabbitMqUserName,
		Password: appConfig.RabbitMqPassword,
		Timeout:  appConfig.RabbitMqTimeout,
		Exchange: logsExchangeName,
	})
	if err != nil {
		panic(err)
	}

	logger := common.NewLogger(rabbitMqClient)

	server := initServer(appConfig, logger)

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

		if err := rabbitMqClient.CloseConnection(); err != nil {
			panic(err)
		}
	}()

	waitGroup.Wait()
}
