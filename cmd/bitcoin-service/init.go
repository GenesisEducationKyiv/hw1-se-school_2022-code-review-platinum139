package main

import (
	"bitcoin-service/config"
	"bitcoin-service/internal/api"
	currencyCtrl "bitcoin-service/internal/currency/controller"
	currency "bitcoin-service/internal/currency/domain"
	"bitcoin-service/internal/currency/infrastructure/currency_provider_creators"
	notificationCtrl "bitcoin-service/internal/notification/controller"
	notification "bitcoin-service/internal/notification/domain"
	subscribersCtrl "bitcoin-service/internal/subscribers/controller"
	subscribers "bitcoin-service/internal/subscribers/domain"
	"bitcoin-service/internal/subscribers/infrastructure"
	"bitcoin-service/pkg/emails"
	"bitcoin-service/pkg/file_storage"
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func initApp() {
	appConfig, err := config.NewAppConfig(".env")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", appConfig.LogLevel)

	fileStorage := file_storage.NewFileStorage(logger, appConfig.StorageFilename)
	subscribersRepo := infrastructure.NewSubscribersFileRepo(*fileStorage)
	subscribersService := subscribers.NewSubscribersService(logger, subscribersRepo)
	subscribersController := subscribersCtrl.NewSubscribersController(logger, appConfig, subscribersService)

	var providerCreator currency.CurrencyProviderCreator
	switch appConfig.CurrencyProvider {
	case string(currency.Coinbase):
		providerCreator = currency_provider_creators.NewCoinbaseProviderCreator(
			logger, appConfig.CachingPeriodMin, appConfig.RateValueBitSize)
	case string(currency.Coingate):
		providerCreator = currency_provider_creators.NewCoingateProviderCreator(
			logger, appConfig.CachingPeriodMin, appConfig.RateValueBitSize)
	default:
		log.Fatalf("Invalid currency provider specified in config.")
	}

	currencyProvider := providerCreator.CreateProvider()
	currencyService := currency.NewCurrencyService(currencyProvider)
	currencyController := currencyCtrl.NewCurrencyController(logger, currencyService)

	mailService := emails.NewEmailService(appConfig)
	notificationService := notification.NewNotificationService(
		logger, appConfig, mailService, currencyService, subscribersService)
	notificationController := notificationCtrl.NewNotificationController(logger, notificationService)

	server := api.NewServer(logger, appConfig, notificationController, currencyController, subscribersController)
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
