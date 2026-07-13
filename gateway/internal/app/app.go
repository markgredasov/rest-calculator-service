package app

import (
	"context"

	clients_computing_core "github.com/markgredasov/rest-calculator-service/internal/clients/computing_core"
	core_logger "github.com/markgredasov/rest-calculator-service/internal/core/logger"
	core_http_server "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/server"
)

type App struct {
	Logger              *core_logger.Logger
	Server              *core_http_server.HTTPServer
	ComputingCoreClient *clients_computing_core.Client
}

func NewApp(
	logger *core_logger.Logger,
	server *core_http_server.HTTPServer,
	computingCoreClient *clients_computing_core.Client,
) *App {
	return &App{
		Logger:              logger,
		Server:              server,
		ComputingCoreClient: computingCoreClient,
	}
}

func (a *App) Run(ctx context.Context) error {
	return a.Server.Run(ctx)
}

func (a *App) Close() error {
	return a.Logger.Close()
}
