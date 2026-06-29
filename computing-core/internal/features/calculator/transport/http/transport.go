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
	Process(ctx context.Context, process domain.Process) (domain.Process, error)
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
			Handler: h.Process,
		},
	}
}
