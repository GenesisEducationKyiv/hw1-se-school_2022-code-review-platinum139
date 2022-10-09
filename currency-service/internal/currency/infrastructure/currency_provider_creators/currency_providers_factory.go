package currency_provider_creators

import (
	currency "currency-service/internal/currency/domain"
	"currency-service/internal/currency/infrastructure/currency_provider_errors"
)

type CurrencyProvidersFactory struct {
	providers map[currency.CurrencyProviderType]currency.CurrencyProviderCreator
}

func (f *CurrencyProvidersFactory) RegisterProvider(
	providerType currency.CurrencyProviderType,
	providerCreator currency.CurrencyProviderCreator,
) {
	f.providers[providerType] = providerCreator
}

func (f *CurrencyProvidersFactory) GetRegistered() []currency.CurrencyProviderType {
	registeredTypes := make([]currency.CurrencyProviderType, 0, len(f.providers))
	for providerType := range f.providers {
		registeredTypes = append(registeredTypes, providerType)
	}
	return registeredTypes
}

func (f *CurrencyProvidersFactory) CreateProvider(
	providerType currency.CurrencyProviderType,
) (currency.Provider, error) {
	providerCreator, registered := f.providers[providerType]
	if !registered {
		return nil, currency_provider_errors.ProviderNotRegisteredError{}
	}

	return providerCreator.CreateProvider(), nil
}

func NewCurrencyProvidersFactory() *CurrencyProvidersFactory {
	return &CurrencyProvidersFactory{
		providers: make(map[currency.CurrencyProviderType]currency.CurrencyProviderCreator),
	}
}
