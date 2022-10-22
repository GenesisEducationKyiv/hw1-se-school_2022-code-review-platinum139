package main

import (
	"subscribers-service/config"
	"subscribers-service/internal/api"
	"subscribers-service/internal/common"
	notificationCtrl "subscribers-service/internal/notification/controller"
	notification "subscribers-service/internal/notification/domain"
	notificationInf "subscribers-service/internal/notification/infrastructure"
	subscribersCtrl "subscribers-service/internal/subscribers/controller"
	subscribers "subscribers-service/internal/subscribers/domain"
	subscribersInf "subscribers-service/internal/subscribers/infrastructure"
	"subscribers-service/pkg/emails"
	"subscribers-service/pkg/file_storage"
)

func initServer(config *config.AppConfig, logger common.Logger) *api.Server {
	fileStorage := file_storage.NewFileStorage(logger, config.StorageFilename)
	subscribersRepo := subscribersInf.NewSubscribersFileRepo(*fileStorage)
	subscribersService := subscribers.NewSubscribersService(logger, subscribersRepo)
	subscribersController := subscribersCtrl.NewSubscribersController(logger, config, subscribersService)

	currencySvcClient := notificationInf.NewCurrencyServiceClientImpl(
		logger, config.CurrencySvcHost, config.CurrencySvcPort, config.RateValueBitSize)

	mailService := emails.NewEmailService(config)
	notificationService := notification.NewNotificationService(
		logger, config, mailService, currencySvcClient, subscribersService)
	notificationController := notificationCtrl.NewNotificationController(logger, notificationService)

	server := api.NewServer(logger, config, notificationController, subscribersController)

	return server
}
