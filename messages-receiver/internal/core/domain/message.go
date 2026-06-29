package domain

type Message struct {
	Message string
}

func NewMessage(message string) Message {
	return Message{
		Message: message,
	}
}
