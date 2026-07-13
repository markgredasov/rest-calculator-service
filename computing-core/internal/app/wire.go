//go:build wireinject
// +build wireinject

package app

import (
	"context"

	"github.com/google/wire"
	core_logger "github.com/markgredasov/rest-calculator-service/internal/core/logger"
	core_http_middleware "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/middleware"
	core_http_server "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/server"
	calculator_service "github.com/markgredasov/rest-calculator-service/internal/features/calculator/service"
	calculator_transport "github.com/markgredasov/rest-calculator-service/internal/features/calculator/transport/http"
)

var ProviderSet = wire.NewSet(
	core_logger.NewConfigMust,
	core_logger.NewLogger,

	provideCalculatorHandler,
	provideAPIRouter,

	core_http_server.NewConfigMust,
	provideMiddlewares,
	provideServer,

	NewApp,
)

func provideCalculatorHandler() *calculator_transport.CalculatorHTTPHandler {
	calculatorService := calculator_service.NewCalculatorService(nil)
	calculatorHTTPHandler := calculator_transport.NewCalculatorHTTPHandler(&calculatorService)

	return calculatorHTTPHandler
}

func provideAPIRouter(
	calculatorHTTPHandler *calculator_transport.CalculatorHTTPHandler,
) *core_http_server.APIVersionRouter {
	apiInternalRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiPrefixInternal)
	apiInternalRouter.RegisterRoutes(calculatorHTTPHandler.Routes()...)

	return apiInternalRouter
}

func provideMiddlewares(
	logger *core_logger.Logger,
) []core_http_middleware.Middleware {
	return []core_http_middleware.Middleware{
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Recovery(),
		core_http_middleware.Trace(),
	}
}

func provideServer(
	config core_http_server.Config,
	logger *core_logger.Logger,
	middlewares []core_http_middleware.Middleware,
	router *core_http_server.APIVersionRouter,
) *core_http_server.HTTPServer {
	server := core_http_server.NewHTTPServer(config, logger, middlewares...)
	server.RegisterAPIRouters(router)

	return server
}

func InitializeApp(ctx context.Context) (*App, error) {
	wire.Build(ProviderSet)
	return nil, nil
}
