package core_http_middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/markgredasov/rest-calculator-service/internal/core/errors"
	core_logger "github.com/markgredasov/rest-calculator-service/internal/core/logger"
	core_http_response "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/response"
	utils_jwt "github.com/markgredasov/rest-calculator-service/internal/utils/jwt"
	"go.uber.org/zap"
)

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		const RequestIDHeader = "x-request-id"
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(RequestIDHeader, requestID)
			w.Header().Set(RequestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(l *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		const RequestIDHeader = "x-request-id"
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIDHeader)

			logger := l.With(
				zap.String("Request-ID", requestID),
				zap.String("Method", r.Method),
				zap.String("URL", r.URL.String()),
			)

			ctx := context.WithValue(r.Context(), "logger", logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Recovery() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/swagger/") {
				next.ServeHTTP(w, r)
				return
			}
			ctx := r.Context()
			logger := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

			defer func() {
				if r := recover(); r != nil {
					responseHandler.PanicResponse(
						r,
						"during handle HTTP request got unexpected panic",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/swagger/") {
				next.ServeHTTP(w, r)
				return
			}
			ctx := r.Context()
			logger := core_logger.FromContext(ctx)

			writer := core_http_response.NewResponseWriter(w)

			start := time.Now().UTC()
			logger.Debug(
				">>> incoming HTTP request",
				zap.Time("time", start),
			)

			next.ServeHTTP(writer, r)

			logger.Debug(
				"<<< done HTTP request",
				zap.Int("status-code", writer.GetStatusCode()),
				zap.Duration("latency", time.Since(start)),
			)
		})
	}
}

func Auth() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

			tokenString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)

			claims, err := utils_jwt.ExtractJWTClaims(tokenString)
			if err != nil {
				err = fmt.Errorf("failed getting claims: %w: %w", err, core_errors.ErrUnauthorized)
				responseHandler.ErrorResponse(
					err,
					"unauthorized",
				)
				return
			}

			r.Header.Set("user_id", claims.UserID)
			r.Header.Set("role", claims.Role)
			next.ServeHTTP(w, r)
		})
	}
}

func RoleRights(checkRole string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

			role := r.Header.Get("role")
			if role != checkRole {
				err := fmt.Errorf("%s rights required: %w", checkRole, core_errors.ErrNoRights)
				responseHandler.ErrorResponse(err, fmt.Sprintf("%s rights required", checkRole))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigins := map[string]struct{}{
				"http://localhost:8080": {},
			}

			origin := r.Header.Get("Origin")

			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
