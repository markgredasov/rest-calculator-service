package utils_http

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	core_errors "github.com/markgredasov/rest-calculator-service/internal/core/errors"
)

func GetUUIDPathValue(r *http.Request, key string) (string, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return "", fmt.Errorf(
			"no key='%s' in path values: %w",
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	if _, err := uuid.Parse(pathValue); err != nil {
		return "", fmt.Errorf(
			"failed to parse path value='%s' by key='%s' to uuid: %w: %w",
			pathValue,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return pathValue, nil
}
