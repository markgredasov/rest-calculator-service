package messages_transport

import (
	"context"
	"net/http"

	"github.com/markgredasov/rest-calculator-service/internal/core/domain"
	core_http_server "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/server"
)

type MessagesHTTPHandler struct {
	MessagesService MessagesService
}

type MessagesService interface {
	SendMessageToCalculator(ctx context.Context, message domain.Message) (domain.Message, error)
}

func NewMessagesHTTPHandler(service MessagesService) *MessagesHTTPHandler {
	return &MessagesHTTPHandler{
		MessagesService: service,
	}
}

func (h *MessagesHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/messages/calculator",
			Handler: h.SendMessageToCalculator,
		},
	}
}
