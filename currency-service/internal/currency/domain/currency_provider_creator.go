package domain

type CurrencyProviderCreator interface {
	CreateProvider() Provider
}
