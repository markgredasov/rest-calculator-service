package messages_service

type MessagesService struct {
	messagesRepository MessagesRepository
}

type MessagesRepository interface{}

func NewMessagesService(repo MessagesRepository) MessagesService {
	return MessagesService{
		messagesRepository: repo,
	}
}
