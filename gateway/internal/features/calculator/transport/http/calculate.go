package calculator_transport

import (
	"net/http"

	"github.com/markgredasov/rest-calculator-service/internal/core/domain"
	core_logger "github.com/markgredasov/rest-calculator-service/internal/core/logger"
	core_http_request "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/request"
	core_http_response "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/response"
	"go.uber.org/zap"
)

type operationParams struct {
	Power int `json:"power" validate:"omitempty,gt=0"`
}

type CalculateDTORequest struct {
	Numbers   []int           `json:"numbers" validate:"required,min=1,dive"`
	Operation string          `json:"operation" validate:"required,oneof=sum multiply average"`
	Params    operationParams `json:"params" validate:"omitempty"`
}

type CalculateDTOResponse struct {
	Status             string  `json:"status"`
	OriginalNumbers    []int   `json:"original_numbers"`
	TransformedNumbers []int   `json:"transformed_numbers"`
	AggregatedResult   float64 `json:"aggregated_result"`
}

func (h *CalculatorHTTPHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("processing calculator handler")
	var request CalculateDTORequest

	// Выполнение валидации на транспортном уровне
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		log.Debug("failed to decode and validate request", zap.Error(err))
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	// Приведение запроса к доменной структуре
	domainRequest := calculateDTOToDomain(request)

	// Валидация на сервисном уровне
	serviceResp, err := h.calculatorService.Calculate(ctx, domainRequest)
	if err != nil {
		log.Debug("failed to calculate result", zap.Any("domain_request", domainRequest))
		responseHandler.ErrorResponse(err, "failed to validate result in service layer")
		return
	}

	// Выполнение запроса к сервису-калькулятору
	result, err := h.computingCoreClient.Calculate(ctx, serviceResp.TransformedNumbers, string(serviceResp.Operation))
	if err != nil {
		log.Error("failed to get result from computing service", zap.Any("domain_request", domainRequest))
		responseHandler.ErrorResponse(err, "failed to get result from computing core")
		return
	}

	// Приведение доменной сущности к структуре ответа
	response := calculateDomainToDTO(serviceResp, result)

	// Отправка ответа
	responseHandler.JSONReponse(response, http.StatusOK)
}

func calculateDTOToDomain(dto CalculateDTORequest) domain.CalculatorRequest {
	var pow int = 1

	if dto.Params.Power != 0 {
		pow = dto.Params.Power
	}

	return domain.NewCalculatorRequestUnitialized(
		dto.Numbers,
		dto.Operation,
		pow,
	)
}

func calculateDomainToDTO(calcDomain domain.CalculatorRequest, result float64) CalculateDTOResponse {
	return CalculateDTOResponse{
		Status:             string(domain.StatusSuccess),
		OriginalNumbers:    calcDomain.OriginalNumbers,
		TransformedNumbers: calcDomain.TransformedNumbers,
		AggregatedResult:   result,
	}
}
