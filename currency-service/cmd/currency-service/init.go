package main

import (
	"context"
	"currency-service/config"
	"currency-service/internal/api"
	"currency-service/internal/currency/controller"
	"currency-service/internal/currency/domain"
	"currency-service/internal/currency/infrastructure/currency_provider_creators"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Run() {
	appConfig, err := config.NewAppConfig(".env")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", appConfig.LogLevel)

	currencyProvidersFactory := initCurrencyProvidersFactory(appConfig)

	currencyProvider, err := initCurrencyProvider(logger, appConfig, currencyProvidersFactory)
	if err != nil {
		log.Fatalf("Failed to create currency provider: %s", err)
	}

	currencyService := domain.NewCurrencyService(currencyProvider)
	currencyController := controller.NewCurrencyController(logger, currencyService)

	server := api.NewServer(logger, appConfig, currencyController)
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
			log.Fatalf("Currency service stopped with error: %s", err)
		}
		log.Println("Currency service stopped gracefully")
	}()

	waitGroup.Wait()
}

func initCurrencyProvidersFactory(config *config.AppConfig) *currency_provider_creators.CurrencyProvidersFactory {
	coinbaseProviderCreator := currency_provider_creators.NewCoinbaseProviderCreator(config.RateValueBitSize)
	coingateProviderCreator := currency_provider_creators.NewCoingateProviderCreator(config.RateValueBitSize)
	kucoinProviderCreator := currency_provider_creators.NewKucoinProviderCreator(config.RateValueBitSize)

	currencyProvidersFactory := currency_provider_creators.NewCurrencyProvidersFactory()
	currencyProvidersFactory.RegisterProvider(domain.CoinbaseProvider, coinbaseProviderCreator)
	currencyProvidersFactory.RegisterProvider(domain.CoingateProvider, coingateProviderCreator)
	currencyProvidersFactory.RegisterProvider(domain.KucoinProvider, kucoinProviderCreator)

	return currencyProvidersFactory
}

func initCurrencyProvider(
	log *log.Logger,
	config *config.AppConfig,
	factory *currency_provider_creators.CurrencyProvidersFactory,
) (domain.Provider, error) {
	providerType := domain.CurrencyProviderType(config.CurrencyProvider)
	provider, err := factory.CreateProvider(providerType)
	if err != nil {
		return nil, err
	}

	prevProvider := provider
	registeredProviderTypes := factory.GetRegistered()
	for _, providerType := range registeredProviderTypes {
		nextProvider, _ := factory.CreateProvider(providerType)
		prevProvider.SetNext(nextProvider)
		prevProvider = nextProvider
	}

	provider = domain.NewCurrencyProviderWithLogger(log, provider)

	provider = domain.NewCachedCurrencyProvider(
		time.Duration(config.CachingPeriodMin)*time.Minute,
		config.RateValueBitSize,
		provider)

	return provider, nil
}
