package domain

type CurrencyService interface {
	GetCurrencyRate(from Currency, to Currency) (float64, error)
}
