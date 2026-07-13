package calculator_transport

import (
	"context"
	"net/http"

	"github.com/markgredasov/rest-calculator-service/internal/core/domain"
	core_http_server "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/server"
)

type CalculatorHTTPHandler struct {
	calculatorService CalculatorService
}

type CalculatorService interface {
	ProcessComputations(ctx context.Context, expression domain.Expression) (domain.Expression, error)
}

func NewCalculatorHTTPHandler(service CalculatorService) *CalculatorHTTPHandler {
	return &CalculatorHTTPHandler{
		calculatorService: service,
	}
}

func (h *CalculatorHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/process",
			Handler: h.ProcessComputations,
		},
	}
}
