package domain

type CurrencyServiceImpl struct {
	currencyProvider Provider
}

func (s *CurrencyServiceImpl) GetCurrencyRate(from Currency, to Currency) (float64, error) {
	return s.currencyProvider.GetCurrencyRate(from, to)
}

func NewCurrencyService(currencyProvider Provider) *CurrencyServiceImpl {
	return &CurrencyServiceImpl{
		currencyProvider: currencyProvider,
	}
}
