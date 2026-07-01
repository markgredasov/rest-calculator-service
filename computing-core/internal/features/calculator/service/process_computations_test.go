package calculator_service_test

import (
	"context"
	"testing"

	"github.com/markgredasov/rest-calculator-service/internal/core/domain"
	core_errors "github.com/markgredasov/rest-calculator-service/internal/core/errors"
	calculator_service "github.com/markgredasov/rest-calculator-service/internal/features/calculator/service"
	"github.com/stretchr/testify/assert"
)

func TestProcessComputations(t *testing.T) {
	service := calculator_service.NewCalculatorService(nil)
	ctx := context.Background()

	tests := []struct {
		name          string
		expression    domain.Expression
		expectedError error
		expectedResp  domain.Expression
	}{
		{
			name: "success: sum operation",
			expression: domain.Expression{
				Numbers:   []int{1, 2, 3, 4, 5},
				Operation: domain.OperationSum,
			},
			expectedError: nil,
			expectedResp: domain.Expression{
				Numbers:   []int{1, 2, 3, 4, 5},
				Operation: domain.OperationSum,
				Result:    75.0,
			},
		},
		{
			name: "success: sum operation with mixed numbers",
			expression: domain.Expression{
				Numbers:   []int{-1, 2, -3, 4, -5},
				Operation: domain.OperationSum,
			},
			expectedError: nil,
			expectedResp: domain.Expression{
				Numbers:   []int{-1, 2, -3, 4, -5},
				Operation: domain.OperationSum,
				Result:    -15.0,
			},
		},
		{
			name: "success: sum operation with zero",
			expression: domain.Expression{
				Numbers:   []int{0, 0, 0},
				Operation: domain.OperationSum,
			},
			expectedError: nil,
			expectedResp: domain.Expression{
				Numbers:   []int{0, 0, 0},
				Operation: domain.OperationSum,
				Result:    0.0,
			},
		},
		{
			name: "success: multiply operation",
			expression: domain.Expression{
				Numbers:   []int{1, 2, 3, 4},
				Operation: domain.OperationMultiply,
			},
			expectedError: nil,
			expectedResp: domain.Expression{
				Numbers:   []int{1, 2, 3, 4},
				Operation: domain.OperationMultiply,
				Result:    6.0,
			},
		},
		{
			name: "success: multiply operation with mixed numbers",
			expression: domain.Expression{
				Numbers:   []int{-1, 2, -3, 4},
				Operation: domain.OperationMultiply,
			},
			expectedError: nil,
			expectedResp: domain.Expression{
				Numbers:   []int{-1, 2, -3, 4},
				Operation: domain.OperationMultiply,
				Result:    6.0,
			},
		},
		{
			name: "success: average operation with positive numbers",
			expression: domain.Expression{
				Numbers:   []int{1, 2, 3, 4, 5},
				Operation: domain.OperationAverage,
			},
			expectedError: nil,
			expectedResp: domain.Expression{
				Numbers:   []int{1, 2, 3, 4, 5},
				Operation: domain.OperationAverage,
				Result:    9.0,
			},
		},
		{
			name: "success: average operation with mixed numbers",
			expression: domain.Expression{
				Numbers:   []int{-1, 2, -3, 4, -5},
				Operation: domain.OperationAverage,
			},
			expectedError: nil,
			expectedResp: domain.Expression{
				Numbers:   []int{-1, 2, -3, 4, -5},
				Operation: domain.OperationAverage,
				Result:    0.36,
			},
		},
		{
			name: "success: average operation with zeros",
			expression: domain.Expression{
				Numbers:   []int{0, 0, 0},
				Operation: domain.OperationAverage,
			},
			expectedError: nil,
			expectedResp: domain.Expression{
				Numbers:   []int{0, 0, 0},
				Operation: domain.OperationAverage,
				Result:    0.0,
			},
		},
		{
			name: "error: empty numbers",
			expression: domain.Expression{
				Numbers:   []int{},
				Operation: domain.OperationSum,
			},
			expectedError: core_errors.ErrInvalidArgument,
			expectedResp:  domain.Expression{},
		},
		{
			name: "error: nil numbers",
			expression: domain.Expression{
				Numbers:   nil,
				Operation: domain.OperationSum,
			},
			expectedError: core_errors.ErrInvalidArgument,
			expectedResp:  domain.Expression{},
		},
		{
			name: "error: invalid operation",
			expression: domain.Expression{
				Numbers:   []int{1, 2, 3},
				Operation: "invalid_op",
			},
			expectedError: core_errors.ErrInvalidArgument,
			expectedResp:  domain.Expression{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.ProcessComputations(ctx, tt.expression)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Equal(t, tt.expectedResp, resp)
			} else {
				assert.NoError(t, err)

				if tt.expectedResp.Result != 0 {
					assert.InDelta(t, tt.expectedResp.Result, resp.Result, 0.0001)
				} else {
					assert.Equal(t, 0.0, resp.Result)
				}

				assert.Equal(t, tt.expectedResp.Numbers, resp.Numbers)
				assert.Equal(t, tt.expectedResp.Operation, resp.Operation)
			}
		})
	}
}
