package calculator_transport

import (
	"net/http"

	"github.com/markgredasov/rest-calculator-service/internal/core/domain"
	core_logger "github.com/markgredasov/rest-calculator-service/internal/core/logger"
	core_http_request "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/request"
	core_http_response "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/response"
	"go.uber.org/zap"
)

type ProcessComputationDTORequest struct {
	Numbers   []int  `json:"numbers" validate:"required,min=1,gte=0"`
	Operation string `json:"operation" validate:"required,oneof=sum multiply average"`
}

type ProcessComputationDTOResponse struct {
	Result float64 `json:"result"`
}

func (h *CalculatorHTTPHandler) ProcessComputations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("processing process computations handler")
	var request ProcessComputationDTORequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		log.Debug("failed to decode and validate request", zap.Error(err))
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	processDomain := processComputationDTOToDomain(request)
	process, err := h.calculatorService.ProcessComputations(ctx, processDomain)
	if err != nil {
		log.Debug("failed to process operation", zap.Any("process_domain", processDomain))
		responseHandler.ErrorResponse(err, "failed to process operation")
		return
	}

	response := processComputationDomainToDTO(process)
	responseHandler.JSONReponse(response, http.StatusOK)
}

func processComputationDTOToDomain(req ProcessComputationDTORequest) domain.Expression {
	return domain.NewExpressionUnitialized(
		req.Numbers,
		req.Operation,
	)
}

func processComputationDomainToDTO(domainExpression domain.Expression) ProcessComputationDTOResponse {
	return ProcessComputationDTOResponse{
		Result: domainExpression.Result,
	}
}
