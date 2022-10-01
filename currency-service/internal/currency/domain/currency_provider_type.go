package domain

type CurrencyProviderType string

const (
	CoinbaseProvider CurrencyProviderType = "coinbase"
	CoingateProvider CurrencyProviderType = "coingate"
	KucoinProvider   CurrencyProviderType = "kucoin"
)
