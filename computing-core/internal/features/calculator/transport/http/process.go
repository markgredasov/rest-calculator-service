package calculator_transport

import (
	"net/http"

	"github.com/markgredasov/rest-calculator-service/internal/core/domain"
	core_logger "github.com/markgredasov/rest-calculator-service/internal/core/logger"
	core_http_request "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/request"
	core_http_response "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/response"
	"go.uber.org/zap"
)

type ProcessDTORequest struct {
	Numbers   []int  `json:"numbers" validate:"required,min=1,gte=0"`
	Operation string `json:"operation" validate:"required,oneof=sum multiply average"`
}

type ProcessDTOResponse struct {
	Result float64 `json:"result"`
}

func (h *CalculatorHTTPHandler) Process(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("processing process handler")
	var request ProcessDTORequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		log.Debug("failed to decode and validate request", zap.Error(err))
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	processDomain := processDTOToDomain(request)
	process, err := h.calculatorService.Process(ctx, processDomain)
	if err != nil {
		log.Debug("failed to process operation", zap.Any("process_domain", processDomain))
		responseHandler.ErrorResponse(err, "failed to process operation")
		return
	}

	response := processDomainToDTO(process)
	responseHandler.JSONReponse(response, http.StatusOK)
}

func processDTOToDomain(req ProcessDTORequest) domain.Process {
	return domain.NewProcessUnitialized(
		req.Numbers,
		req.Operation,
	)
}

func processDomainToDTO(domainProcess domain.Process) ProcessDTOResponse {
	return ProcessDTOResponse{
		Result: domainProcess.Result,
	}
}
