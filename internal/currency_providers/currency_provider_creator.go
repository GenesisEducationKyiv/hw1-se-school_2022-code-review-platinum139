package currency_providers

import "bitcoin-service/pkg/currency"

type CurrencyProviderCreator interface {
	CreateProvider() currency.Provider
}
