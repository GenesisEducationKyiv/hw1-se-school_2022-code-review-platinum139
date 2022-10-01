package domain

type CurrencyServiceClient interface {
	GetRate() (float64, error)
}
