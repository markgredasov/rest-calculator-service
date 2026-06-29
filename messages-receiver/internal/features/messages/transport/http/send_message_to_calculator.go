package messages_transport

import (
	"net/http"

	"github.com/markgredasov/rest-calculator-service/internal/core/domain"
	core_logger "github.com/markgredasov/rest-calculator-service/internal/core/logger"
	core_http_request "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/request"
	core_http_response "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/response"
)

type SendMessageDTORequest struct {
	Message string `json:"message" validate:"required,min=1"`
}

type SendMessageDTOResponse struct {
	Result string `json:"result"`
}

func (h *MessagesHTTPHandler) SendMessageToCalculator(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("processing send message to calculator handler")
	var request SendMessageDTORequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	domain := sendMessageDTOToDomain(request)
	message, err := h.MessagesService.SendMessageToCalculator(ctx, domain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to send message")
		return
	}

	response := sendMessageDomainToDTO(message)
	responseHandler.JSONReponse(response, http.StatusOK)
}

func sendMessageDTOToDomain(dto SendMessageDTORequest) domain.Message {
	return domain.NewMessage(
		dto.Message,
	)
}

func sendMessageDomainToDTO(msg domain.Message) SendMessageDTOResponse {
	return SendMessageDTOResponse{
		Result: msg.Message,
	}
}
