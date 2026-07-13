//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"time"

	"github.com/google/wire"
	clients_computing_core "github.com/markgredasov/rest-calculator-service/internal/clients/computing_core"
	core_logger "github.com/markgredasov/rest-calculator-service/internal/core/logger"
	core_http_middleware "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/middleware"
	core_http_server "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/server"
	calculator_service "github.com/markgredasov/rest-calculator-service/internal/features/calculator/service"
	calculator_transport "github.com/markgredasov/rest-calculator-service/internal/features/calculator/transport/http"
)

var ProviderSet = wire.NewSet(
	core_logger.NewConfigMust,
	core_logger.NewLogger,

	provideComputingCoreClient,
	provideCalculatorService,

	wire.Bind(new(calculator_transport.ComputingCoreClient), new(*clients_computing_core.Client)),
	wire.Bind(new(calculator_transport.CalculatorService), new(*calculator_service.CalculatorService)),

	calculator_transport.NewCalculatorHTTPHandler,

	provideAPIRouter,
	provideServerConfig,
	provideMiddlewares,
	provideServer,

	NewApp,
)

func provideComputingCoreClient() *clients_computing_core.Client {
	return clients_computing_core.NewClient(3 * time.Second)
}

func provideCalculatorService() *calculator_service.CalculatorService {
	service := calculator_service.NewCalculatorService(nil)
	return &service
}

func provideServerConfig() core_http_server.Config {
	return core_http_server.NewConfigMust()
}

func provideAPIRouter(
	calculatorHandler *calculator_transport.CalculatorHTTPHandler,
) *core_http_server.APIVersionRouter {
	router := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	router.RegisterRoutes(calculatorHandler.Routes()...)
	return router
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
