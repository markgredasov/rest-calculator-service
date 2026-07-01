package clients_computing_core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	clients_computing_core_models "github.com/markgredasov/rest-calculator-service/internal/clients/computing_core/models"
	core_errors "github.com/markgredasov/rest-calculator-service/internal/core/errors"
)

func (c *Client) Calculate(ctx context.Context, numbers []int, operation string) (float64, error) {
	const RequestIDHeader = "x-request-id"

	reqStruct := clients_computing_core_models.NewRequest(
		numbers,
		operation,
	)

	reqJSON, err := json.Marshal(reqStruct)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	requestID := ctx.Value(RequestIDHeader).(string)
	if requestID == "" {
		requestID = uuid.NewString()
		ctx = context.WithValue(ctx, "Request-ID", requestID)
	}

	httpCall := func() (interface{}, error) {
		ctxForRequest, cancel := context.WithTimeout(ctx, c.Timeout)
		defer cancel()

		req, err := http.NewRequestWithContext(ctxForRequest, http.MethodPost, c.baseURL, bytes.NewBuffer(reqJSON))
		if err != nil {
			return 0, fmt.Errorf("failed to initialize request with context: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(RequestIDHeader, requestID)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return 0, fmt.Errorf("failed to create request: %w: %w", err, core_errors.ErrServiceUnavailable)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return 0, fmt.Errorf("computing service returned error: %w: %w", err, core_errors.ErrBadGateway)
		}

		var responseStruct clients_computing_core_models.Response
		if err := json.NewDecoder(resp.Body).Decode(&responseStruct); err != nil {
			return 0, fmt.Errorf("failed to decode response: %w", err)
		}

		return responseStruct.Result, nil
	}

	result, err := c.circuitBreaker.Execute(httpCall)
	if err != nil {
		return 0, fmt.Errorf("failed to execute circuit breaker: %w: %w", err, core_errors.ErrServiceUnavailable)
	}

	return result.(float64), nil
}
