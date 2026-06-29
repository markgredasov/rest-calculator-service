package domain

import "fmt"

type Operation string

const (
	OperationSum      = Operation("sum")
	OperationMultiply = Operation("multiply")
	OperationAverage  = Operation("average")
)

type Process struct {
	Numbers   []int     `json:"numbers"`
	Operation Operation `json:"operation"`
	Result    float64   `json:"result"`
}

func NewProcess(nums []int, op string, result float64) Process {
	return Process{
		Numbers:   nums,
		Operation: Operation(op),
		Result:    result,
	}
}

func NewProcessUnitialized(nums []int, op string) Process {
	return Process{
		Numbers:   nums,
		Operation: Operation(op),
		Result:    UnitializedResult,
	}
}

func (p *Process) Validate() error {

	if len(p.Numbers) == 0 {
		return fmt.Errorf("empty numbers array")
	}

	switch p.Operation {
	case OperationSum, OperationMultiply, OperationAverage:
		break
	default:
		return fmt.Errorf("invalid operation value")
	}

	return nil
}
