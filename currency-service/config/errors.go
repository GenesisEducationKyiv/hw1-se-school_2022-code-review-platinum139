package config

type LoadError struct {
	Message string
}

func (e LoadError) Error() string {
	return e.Message
}
