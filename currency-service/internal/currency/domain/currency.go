package domain

type Currency string

func (c Currency) ToString() string {
	return string(c)
}

const (
	Uah Currency = "UAH"
	Btc Currency = "BTC"
)
