package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/markgredasov/rest-calculator-service/internal/core/errors"
	core_logger "github.com/markgredasov/rest-calculator-service/internal/core/logger"
	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	log *core_logger.Logger
	w   http.ResponseWriter
}

func NewHTTPResponseHandler(l *core_logger.Logger, w http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: l,
		w:   w,
	}
}

func (h *HTTPResponseHandler) JSONReponse(responseBody any, statusCode int) {
	h.w.WriteHeader(statusCode)

	if err := json.NewEncoder(h.w).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrNotFound):
		err = core_errors.ErrNotFound
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug
	case errors.Is(err, core_errors.ErrConflict):
		err = core_errors.ErrConflict
		statusCode = http.StatusConflict
		logFunc = h.log.Warn
	case errors.Is(err, core_errors.ErrInvalidArgument):
		err = core_errors.ErrInvalidArgument
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn
	case errors.Is(err, core_errors.ErrUnauthorized):
		err = core_errors.ErrUnauthorized
		statusCode = http.StatusUnauthorized
		logFunc = h.log.Debug
	case errors.Is(err, core_errors.ErrNoRights):
		err = core_errors.ErrNoRights
		statusCode = http.StatusForbidden
		logFunc = h.log.Debug
	case errors.Is(err, core_errors.ErrBadGateway):
		err = core_errors.ErrBadGateway
		statusCode = http.StatusBadGateway
		logFunc = h.log.Error
	case errors.Is(err, core_errors.ErrServiceUnavailable):
		err = core_errors.ErrServiceUnavailable
		statusCode = http.StatusServiceUnavailable
		logFunc = h.log.Error
	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) errorResponse(statusCode int, err error, msg string) {
	errorResponse := ErrorResponse{
		Error:   err.Error(),
		Message: msg,
	}

	h.JSONReponse(errorResponse, statusCode)
}
