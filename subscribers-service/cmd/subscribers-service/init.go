package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"subscribers-service/config"
	"subscribers-service/internal/api"
	notificationCtrl "subscribers-service/internal/notification/controller"
	notification "subscribers-service/internal/notification/domain"
	notificationInf "subscribers-service/internal/notification/infrastructure"
	subscribersCtrl "subscribers-service/internal/subscribers/controller"
	subscribers "subscribers-service/internal/subscribers/domain"
	subscribersInf "subscribers-service/internal/subscribers/infrastructure"
	"subscribers-service/pkg/emails"
	"subscribers-service/pkg/file_storage"
	"sync"
	"syscall"
)

func Run() {
	appConfig, err := config.NewAppConfig(".env")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", appConfig.LogLevel)

	fileStorage := file_storage.NewFileStorage(logger, appConfig.StorageFilename)
	subscribersRepo := subscribersInf.NewSubscribersFileRepo(*fileStorage)
	subscribersService := subscribers.NewSubscribersService(logger, subscribersRepo)
	subscribersController := subscribersCtrl.NewSubscribersController(logger, appConfig, subscribersService)

	currencySvcClient := notificationInf.NewCurrencyServiceClientImpl(
		logger, appConfig.CurrencySvcHost, appConfig.CurrencySvcPort, appConfig.RateValueBitSize)

	mailService := emails.NewEmailService(appConfig)
	notificationService := notification.NewNotificationService(
		logger, appConfig, mailService, currencySvcClient, subscribersService)
	notificationController := notificationCtrl.NewNotificationController(logger, notificationService)

	server := api.NewServer(logger, appConfig, notificationController, subscribersController)
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
