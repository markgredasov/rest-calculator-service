package app

import (
	"context"

	core_logger "github.com/markgredasov/rest-calculator-service/internal/core/logger"
	core_http_server "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/server"
)

type App struct {
	Logger *core_logger.Logger
	Server *core_http_server.HTTPServer
}

func NewApp(
	logger *core_logger.Logger,
	server *core_http_server.HTTPServer,
) *App {
	return &App{
		Logger: logger,
		Server: server,
	}
}

func (a *App) Run(ctx context.Context) error {
	return a.Server.Run(ctx)
}

func (a *App) Close() error {
	return a.Logger.Close()
}
