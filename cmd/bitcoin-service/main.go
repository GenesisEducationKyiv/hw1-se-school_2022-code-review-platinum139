package main

import (
	"bitcoin-service/config"
	"bitcoin-service/internal/api"
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

	mailService := emails.NewEmailService(appConfig)
	store := storage.NewFilestore(logger, appConfig.StorageFilename)
	subscribersService := subscribers.NewSubscribersService(logger, store, mailService)
	currencyService := currency.NewCurrencyService(appConfig.CurrencyServiceBaseURL)

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
