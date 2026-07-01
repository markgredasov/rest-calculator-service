package calculator_service_test

import (
	"context"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/markgredasov/rest-calculator-service/internal/core/domain"
	core_errors "github.com/markgredasov/rest-calculator-service/internal/core/errors"
	calculator_service "github.com/markgredasov/rest-calculator-service/internal/features/calculator/service"
)

func TestCalculate(t *testing.T) {
	service := calculator_service.NewCalculatorService(nil)
	ctx := context.Background()

	tests := []struct {
		name          string
		request       domain.CalculatorRequest
		expectedError error
		expectedResp  domain.CalculatorRequest
	}{
		{
			name: "success: valid request with addition",
			request: domain.CalculatorRequest{
				OriginalNumbers:    []int{1, 2, 3},
				TransformedNumbers: []int{1, 4, 9},
				Operation:          "sum",
				Power:              2,
			},
			expectedError: nil,
			expectedResp: domain.CalculatorRequest{
				OriginalNumbers:    []int{1, 2, 3},
				TransformedNumbers: []int{1, 4, 9},
				Operation:          "sum",
				Power:              2,
			},
		},
		{
			name: "error: empty numbers",
			request: domain.CalculatorRequest{
				OriginalNumbers:    []int{},
				TransformedNumbers: []int{},
				Operation:          "sum",
				Power:              2,
			},
			expectedError: core_errors.ErrInvalidArgument,
			expectedResp:  domain.CalculatorRequest{},
		},
		{
			name: "error: nil numbers",
			request: domain.CalculatorRequest{
				OriginalNumbers:    nil,
				TransformedNumbers: nil,
				Operation:          "sum",
				Power:              2,
			},
			expectedError: core_errors.ErrInvalidArgument,
			expectedResp:  domain.CalculatorRequest{},
		},
		{
			name: "error: invalid operation",
			request: domain.CalculatorRequest{
				OriginalNumbers:    []int{1, 2, 3},
				TransformedNumbers: []int{1, 4, 9},
				Operation:          "invalid operation",
				Power:              2,
			},
			expectedError: core_errors.ErrInvalidArgument,
			expectedResp:  domain.CalculatorRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.Calculate(ctx, tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Equal(t, tt.expectedResp, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}
