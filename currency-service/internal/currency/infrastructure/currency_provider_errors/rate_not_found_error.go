package currency_provider_errors

type RateNotFoundError struct{}

func (e RateNotFoundError) Error() string {
	return "rate not found"
}
