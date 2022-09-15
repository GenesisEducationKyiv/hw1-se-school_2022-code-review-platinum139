package main

import (
	"bitcoin-service/config"
	"bitcoin-service/internal/api"
	"bitcoin-service/internal/currency_providers"
	"bitcoin-service/internal/repository"
	"bitcoin-service/internal/subscribers"
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

	var providerCreator currency_providers.CurrencyProviderCreator
	switch appConfig.CurrencyProvider {
	case string(currency_providers.Coinbase):
		providerCreator = currency_providers.NewCoinbaseProviderCreator(logger, appConfig.CachingPeriodMin)
	case string(currency_providers.Coingate):
		providerCreator = currency_providers.NewCoingateProviderCreator(logger, appConfig.CachingPeriodMin)
	default:
		log.Fatalf("Invalid currency provider specified in config.")
	}

	currencyService := providerCreator.CreateProvider()

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
