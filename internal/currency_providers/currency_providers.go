package currency_providers

type CurrencyProvider string

const (
	Coinbase CurrencyProvider = "coinbase"
	Coingate CurrencyProvider = "coingate"
	Kucoin   CurrencyProvider = "kucoin"
)
