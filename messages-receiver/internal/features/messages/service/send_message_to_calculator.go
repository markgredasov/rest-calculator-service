package messages_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/markgredasov/rest-calculator-service/internal/core/domain"
	core_errors "github.com/markgredasov/rest-calculator-service/internal/core/errors"
)

func (s *MessagesService) SendMessageToCalculator(ctx context.Context, message domain.Message) (domain.Message, error) {
	msg := message.Message

	msg = strings.TrimSpace(msg)

	if len(msg) < 1 {
		return domain.Message{}, fmt.Errorf("empty message: %w", core_errors.ErrInvalidArgument)
	}

	for _, v := range msg {
		if !isValidSymbol(v) {
			return domain.Message{}, fmt.Errorf("invalid character '%c' in expression: %w", v, core_errors.ErrInvalidArgument)
		}
	}

	calcRequest := domain.CalculatorRequest{
		Expression: msg,
	}

	jsonData, err := json.Marshal(calcRequest)
	if err != nil {
		return domain.Message{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := "http://localhost:8081/api/v1/calculate"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return domain.Message{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return domain.Message{}, fmt.Errorf("failed to send request to calculator: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.Message{}, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return domain.Message{}, fmt.Errorf("calculator service returned error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var calcResponse domain.CalculatorResponse
	if err := json.Unmarshal(body, &calcResponse); err != nil {
		return domain.Message{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	resultMessage := domain.Message{
		Message: strconv.FormatFloat(calcResponse.Result, 'f', 2, 64),
	}

	return resultMessage, nil
}

func isValidSymbol(ch rune) bool {
	switch {
	case ch >= '0' && ch <= '9':
		return true
	case ch == '+' || ch == '-' || ch == '*' || ch == '/':
		return true
	case ch == '(' || ch == ')' || ch == '.' || ch == ' ':
		return true
	default:
		return false
	}
}
