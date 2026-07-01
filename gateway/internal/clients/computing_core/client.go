package clients_computing_core

import (
	"fmt"
	"net/http"
	"time"

	utils_env "github.com/markgredasov/rest-calculator-service/internal/utils/env"
	"github.com/sony/gobreaker"
)

const (
	host     = "HTTP_HOST"
	addr     = "HTTP_ADDR_COMPUTING_CORE"
	endpoint = "CALCULATE_URL"
)

type Client struct {
	Timeout time.Duration

	baseURL        string
	httpClient     *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

func NewClient(timeout time.Duration) *Client {
	baseURL, err := getBaseURL()
	if err != nil {
		err = fmt.Errorf("failed to get base URL: %w", err)
		panic(err)
	}

	cb := initCircuitBreaker()

	return &Client{
		Timeout: timeout,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		circuitBreaker: cb,
	}
}

func initCircuitBreaker() *gobreaker.CircuitBreaker {
	cbSettings := gobreaker.Settings{
		Name:        "computing-core-client",
		MaxRequests: 3,
		Timeout:     15 * time.Second,

		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failurePercent := float64(counts.TotalFailures) / float64(counts.Requests)
			return failurePercent > 0.5 && counts.Requests > 10
		},
	}

	cb := gobreaker.NewCircuitBreaker(cbSettings)

	return cb
}

func getBaseURL() (string, error) {
	environments := []string{
		host,
		addr,
		endpoint,
	}

	envs, err := utils_env.GetEnvironments(environments...)
	if err != nil {
		return "", fmt.Errorf("failed to get environments: %w", err)
	}

	return fmt.Sprintf(
		"http://%s%s/%s",
		envs[host],
		envs[addr],
		envs[endpoint],
	), nil
}
