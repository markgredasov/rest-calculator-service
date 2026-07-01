package calculator_service

import (
	"context"
	"fmt"

	"github.com/markgredasov/rest-calculator-service/internal/core/domain"
	core_errors "github.com/markgredasov/rest-calculator-service/internal/core/errors"
)

func (s *CalculatorService) Calculate(ctx context.Context,
	calculatorRequest domain.CalculatorRequest) (domain.CalculatorRequest, error) {
	if err := calculatorRequest.Validate(); err != nil {
		return domain.CalculatorRequest{},
			fmt.Errorf("failed to validate calculator request: %w", core_errors.ErrInvalidArgument)
	}

	return calculatorRequest, nil
}
