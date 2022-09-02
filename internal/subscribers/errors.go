package subscribers

import "strings"

const separator = ","

type SendMailError struct {
	Subscribers []string
}

func (e SendMailError) Error() string {
	emails := strings.Join(e.Subscribers, separator)
	return "Sending email failed for: " + emails
}
