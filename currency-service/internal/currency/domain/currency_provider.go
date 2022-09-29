package domain

type Provider interface {
	GetCurrencyRate(from Currency, to Currency) (float64, error)
	SetNext(next Provider)
}
