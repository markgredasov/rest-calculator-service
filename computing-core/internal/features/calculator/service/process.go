package calculator_service

import (
	"context"
	"fmt"
	"math"

	"github.com/markgredasov/rest-calculator-service/internal/core/domain"
	core_errors "github.com/markgredasov/rest-calculator-service/internal/core/errors"
)

func (s *CalculatorService) Process(ctx context.Context, process domain.Process) (domain.Process, error) {
	if err := process.Validate(); err != nil {
		return domain.Process{}, fmt.Errorf("failed to validate: %w", core_errors.ErrInvalidArgument)
	}

	var result float64
	length := float64(len(process.Numbers))

	switch process.Operation {
	case domain.OperationSum:
		var sum float64
		for _, num := range process.Numbers {
			sum += float64(num)
		}

		result = sum * length
	case domain.OperationMultiply:
		var mul float64 = 1
		for _, num := range process.Numbers {
			mul *= float64(num)
		}

		result = mul / length
	case domain.OperationAverage:
		var sum float64
		for _, num := range process.Numbers {
			sum += float64(num)
		}

		result = sum / length
		result = math.Pow(result, 2)
	default:
		return domain.Process{}, fmt.Errorf("invalid operation: %w", core_errors.ErrInvalidArgument)
	}

	resultProcess := domain.NewProcess(
		process.Numbers,
		string(process.Operation),
		result,
	)

	return resultProcess, nil
}
