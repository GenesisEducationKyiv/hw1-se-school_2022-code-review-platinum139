package main

import (
	"bitcoin-service/config"
	"bitcoin-service/internal/api"
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
	appConfig, err := config.NewAppConfig()
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", appConfig.LogLevel)

	mailService := emails.NewEmailService(appConfig)
	store := storage.NewFilestore(logger, appConfig.StorageFilename)
	service := subscribers.NewService(logger, store, mailService)

	server := api.NewServer(logger, appConfig, service)
	server.RegisterRoutes()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Run()
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		<-signals
		defer wg.Done()

		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("Server stopped with error: %s", err)
		}
		log.Println("Server stopped gracefully")
	}()

	wg.Wait()
}
