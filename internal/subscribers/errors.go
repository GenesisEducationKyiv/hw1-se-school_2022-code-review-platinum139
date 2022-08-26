package subscribers

import "strings"

type SendMailError struct {
	Subscribers []string
}

func (e SendMailError) Error() string {
	emails := strings.Join(e.Subscribers, ",")
	return "Sending email failed for: " + emails
}
