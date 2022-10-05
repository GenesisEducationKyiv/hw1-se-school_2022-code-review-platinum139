package currency_provider_errors

type ProviderNotRegisteredError struct{}

func (e ProviderNotRegisteredError) Error() string {
	return "provider is not registered"
}
