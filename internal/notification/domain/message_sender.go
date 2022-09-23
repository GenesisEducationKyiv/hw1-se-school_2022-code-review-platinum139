package domain

type MessageSender interface {
	Send(receiver string, subject string, text string) error
}
