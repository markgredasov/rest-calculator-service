package domain

import (
	"fmt"
	"math"
)

type Operation string

const (
	OperationSum      = Operation("sum")
	OperationMultiply = Operation("multiply")
	OperationAverage  = Operation("average")
)

type Status string

const (
	StatusSuccess = Status("success")
	StatusFailed  = Status("failed")
)

type CalculatorRequest struct {
	Status             Status    `json:"status"`
	OriginalNumbers    []int     `json:"original_numbers"`
	TransformedNumbers []int     `json:"transformed_numbers"`
	AggregatedResult   float64   `json:"aggregated_result"`
	Operation          Operation `json:"operation"`
	Power              int       `json:"power"`
}

func NewCalculatorRequestUnitialized(nums []int, operation string, power int) CalculatorRequest {
	transformedNumbers := make([]int, len(nums))

	for i, value := range nums {
		transformedNumbers[i] = int(math.Pow(float64(value), float64(power)))
	}

	return CalculatorRequest{
		Status:             Status(UnitializedStatus),
		OriginalNumbers:    nums,
		TransformedNumbers: transformedNumbers,
		AggregatedResult:   UnitializedAggregatedResult,
		Operation:          Operation(operation),
		Power:              power,
	}
}

func NewCalculatorRequest(status Status, originalNumbers []int, transformedNumbers []int,
	aggregatedResult float64, operation Operation, power int) CalculatorRequest {
	return CalculatorRequest{
		Status:             status,
		OriginalNumbers:    originalNumbers,
		TransformedNumbers: transformedNumbers,
		AggregatedResult:   aggregatedResult,
		Operation:          operation,
		Power:              power,
	}
}

func (c *CalculatorRequest) Validate() error {
	if len(c.OriginalNumbers) == 0 {
		return fmt.Errorf("empty original numbers array: %d", len(c.OriginalNumbers))
	}

	switch c.Operation {
	case OperationSum, OperationMultiply, OperationAverage:
		break
	default:
		return fmt.Errorf("invalid operation: %s", c.Operation)
	}

	if c.Power <= 1 {
		return fmt.Errorf("invalid power: %d", c.Power)
	}

	return nil
}
