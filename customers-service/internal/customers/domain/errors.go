package domain

type EmailAlreadyExistsError struct{}

func (EmailAlreadyExistsError) Error() string {
	return "email already exists"
}
