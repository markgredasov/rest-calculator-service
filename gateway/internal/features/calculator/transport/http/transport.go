package calculator_transport

import (
	"context"
	"net/http"

	"github.com/markgredasov/rest-calculator-service/internal/core/domain"
	core_http_server "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/server"
)

type CalculatorHTTPHandler struct {
	calculatorService   CalculatorService
	computingCoreClient ComputingCoreClient
}

type CalculatorService interface {
	Calculate(ctx context.Context, calculatorRequest domain.CalculatorRequest) (domain.CalculatorRequest, error)
}

type ComputingCoreClient interface {
	Calculate(ctx context.Context, numbers []int, operation string) (float64, error)
}

func NewCalculatorHTTPHandler(service CalculatorService, client ComputingCoreClient) *CalculatorHTTPHandler {
	return &CalculatorHTTPHandler{
		calculatorService:   service,
		computingCoreClient: client,
	}
}

func (h *CalculatorHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/calculate",
			Handler: h.Calculate,
		},
	}
}
