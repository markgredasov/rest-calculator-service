package calculator_service

import (
	"context"
	"fmt"
	"math"

	"github.com/markgredasov/rest-calculator-service/internal/core/domain"
	core_errors "github.com/markgredasov/rest-calculator-service/internal/core/errors"
)

func (s *CalculatorService) ProcessComputations(ctx context.Context, expression domain.Expression) (domain.Expression, error) {
	if err := expression.Validate(); err != nil {
		return domain.Expression{}, fmt.Errorf("failed to validate: %w", core_errors.ErrInvalidArgument)
	}

	var result float64
	length := float64(len(expression.Numbers))

	switch expression.Operation {
	case domain.OperationSum:
		var sum float64
		for _, num := range expression.Numbers {
			sum += float64(num)
		}

		result = sum * length
	case domain.OperationMultiply:
		var mul float64 = 1
		for _, num := range expression.Numbers {
			mul *= float64(num)
		}

		result = mul / length
	case domain.OperationAverage:
		var sum float64
		for _, num := range expression.Numbers {
			sum += float64(num)
		}

		result = sum / length
		result = math.Pow(result, 2)
	default:
		return domain.Expression{}, fmt.Errorf("invalid operation: %w", core_errors.ErrInvalidArgument)
	}

	resultProcess := domain.NewExpression(
		expression.Numbers,
		string(expression.Operation),
		result,
	)

	return resultProcess, nil
}
