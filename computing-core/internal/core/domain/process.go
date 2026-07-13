package domain

import "fmt"

type Operation string

const (
	OperationSum      = Operation("sum")
	OperationMultiply = Operation("multiply")
	OperationAverage  = Operation("average")
)

type Expression struct {
	Numbers   []int     `json:"numbers"`
	Operation Operation `json:"operation"`
	Result    float64   `json:"result"`
}

func NewExpression(nums []int, op string, result float64) Expression {
	return Expression{
		Numbers:   nums,
		Operation: Operation(op),
		Result:    result,
	}
}

func NewExpressionUnitialized(nums []int, op string) Expression {
	return Expression{
		Numbers:   nums,
		Operation: Operation(op),
		Result:    UnitializedResult,
	}
}

func (e *Expression) Validate() error {

	if len(e.Numbers) == 0 {
		return fmt.Errorf("empty numbers array")
	}

	switch e.Operation {
	case OperationSum, OperationMultiply, OperationAverage:
		break
	default:
		return fmt.Errorf("invalid operation value")
	}

	return nil
}
