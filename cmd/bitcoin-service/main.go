package main

import (
	"bitcoin-service/config"
	"bitcoin-service/internal/api"
	"bitcoin-service/internal/repository"
	"bitcoin-service/internal/subscribers"
	"bitcoin-service/pkg/currency"
	"bitcoin-service/pkg/emails"
	"bitcoin-service/pkg/storage"
	"context"
	"log"
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

	logger := log.New(os.Stdout, "", appConfig.LogLevel)

	fileStorage := storage.NewFilestore(logger, appConfig.StorageFilename)
	mailService := emails.NewEmailService(appConfig)
	subscribersRepo := repository.NewSubscribersFileRepo(fileStorage)
	subscribersService := subscribers.NewSubscribersService(logger, subscribersRepo, mailService)
	currencyService := currency.NewCurrencyService()

	server := api.NewServer(logger, appConfig, subscribersService, currencyService)
	server.RegisterRoutes()

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		server.Run()
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	waitGroup.Add(1)
	go func() {
		<-signals
		defer waitGroup.Done()

		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("Server stopped with error: %s", err)
		}
		log.Println("Server stopped gracefully")
	}()

	waitGroup.Wait()
}
