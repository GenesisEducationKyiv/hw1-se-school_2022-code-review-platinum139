package domain

type CurrencyProviderType string

const (
	Coinbase CurrencyProviderType = "coinbase"
	Coingate CurrencyProviderType = "coingate"
	Kucoin   CurrencyProviderType = "kucoin"
)
