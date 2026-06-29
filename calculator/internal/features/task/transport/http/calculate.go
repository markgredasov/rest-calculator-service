package tasks_transport

import (
	"net/http"

	"github.com/markgredasov/rest-calculator-service/internal/core/domain"
	core_logger "github.com/markgredasov/rest-calculator-service/internal/core/logger"
	core_http_request "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/request"
	core_http_response "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/response"
)

type CalculateDTORequest struct {
	Expression string `json:"expression" validate:"required"`
}

type CalculateDTOResponse struct {
	Result float64 `json:"result"`
}

func (h *TasksHTTPHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("processing calculate handler")
	var request CalculateDTORequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	domain := calculateDTOToDomain(request)
	task, err := h.tasksService.Calculate(ctx, domain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get calculate")
		return
	}

	resp := calculateDomainToDTO(task)
	responseHandler.JSONReponse(resp, http.StatusOK)
}

func calculateDTOToDomain(dto CalculateDTORequest) domain.Task {
	return domain.Task{
		Expression: dto.Expression,
	}
}

func calculateDomainToDTO(domain domain.Task) CalculateDTOResponse {
	return CalculateDTOResponse{
		Result: domain.Result,
	}
}
