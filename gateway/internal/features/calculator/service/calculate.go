package calculator_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/markgredasov/rest-calculator-service/internal/core/domain"
	core_errors "github.com/markgredasov/rest-calculator-service/internal/core/errors"
)

type computingCoreResponse struct {
	Result float64 `json:"result"`
}

func (s *CalculatorService) Calculate(
	ctx context.Context,
	calculatorRequest domain.CalculatorRequest,
) (domain.CalculatorRequest, error) {
	if err := calculatorRequest.Validate(); err != nil {
		return domain.CalculatorRequest{},
			fmt.Errorf("failed to validate calculator request: %w", core_errors.ErrInvalidArgument)
	}

	jsonData, err := calculatorRequestToJSONRequest(calculatorRequest)
	if err != nil {
		return domain.CalculatorRequest{},
			fmt.Errorf("failed to marshall request: %w", err)
	}

	url, err := getCamputingCoreURL()
	if err != nil {
		return domain.CalculatorRequest{},
			fmt.Errorf("failed to get computing core URL: %w", err)
	}

	ctxForRequest, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxForRequest, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return domain.CalculatorRequest{},
			fmt.Errorf("failed to initialize request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return domain.CalculatorRequest{},
			fmt.Errorf("failed to create request: %w: %w", err, core_errors.ErrServiceUnavailable)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.CalculatorRequest{},
			fmt.Errorf("computing service returned error: %w: %w", err, core_errors.ErrBadGateway)
	}

	var resultResponse computingCoreResponse

	if err := json.NewDecoder(resp.Body).Decode(&resultResponse); err != nil {
		return domain.CalculatorRequest{},
			fmt.Errorf("failed to decode response: %w", err)
	}

	calculatorResponse := domain.CalculatorRequest{
		Status:             domain.StatusSuccess,
		OriginalNumbers:    calculatorRequest.OriginalNumbers,
		TransformedNumbers: calculatorRequest.TransformedNumbers,
		Operation:          calculatorRequest.Operation,
		Power:              calculatorRequest.Power,
		AggregatedResult:   resultResponse.Result,
	}

	return calculatorResponse, nil
}

func calculatorRequestToJSONRequest(calculatorRequest domain.CalculatorRequest) ([]byte, error) {
	data := map[string]interface{}{
		"numbers":   calculatorRequest.TransformedNumbers,
		"operation": calculatorRequest.Operation,
	}

	return json.Marshal(data)
}

func getCamputingCoreURL() (string, error) {
	host := os.Getenv("HTTP_HOST")
	if host == "" {
		return "", fmt.Errorf("failed to get HTTP_HOST")
	}

	addr := os.Getenv("HTTP_ADDR_COMPUTING_CORE")
	if addr == "" {
		return "", fmt.Errorf("failed to get HTTP_ADDR_COMPUTING_CORE")
	}

	calculateURL := os.Getenv("CALCULATE_URL")
	if calculateURL == "" {
		return "", fmt.Errorf("failed to get CALCULATE_URL")
	}

	return fmt.Sprintf(
		"http://%s%s/%s",
		host,
		addr,
		calculateURL,
	), nil
}
