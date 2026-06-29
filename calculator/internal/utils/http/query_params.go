package utils_http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	core_errors "github.com/markgredasov/rest-calculator-service/internal/core/errors"
)

func GetDateQueryParam(r *http.Request, key string) (*time.Time, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, fmt.Errorf(
			"no key='%s' in query parameters: %w",
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	date, err := time.Parse("2006-01-02", param)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse date='%s' by key='%s': %w: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &date, nil
}

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to convert param='%s' by key='%s' in query parameters: %w",
			param,
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	return &val, nil
}
