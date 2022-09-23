package domain

import "strings"

const separator = ","

type SendMessageError struct {
	FailedSubscribers []string
}

func (e SendMessageError) Error() string {
	emails := strings.Join(e.FailedSubscribers, separator)
	return "Sending email failed for: " + emails
}
