package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/markgredasov/rest-calculator-service/internal/core/logger"
	core_http_middleware "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/middleware"
	core_http_server "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/server"
	messages_service "github.com/markgredasov/rest-calculator-service/internal/features/messages/service"
	messages_transport "github.com/markgredasov/rest-calculator-service/internal/features/messages/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	loggerConfig := core_logger.NewConfigMust()
	logger, err := core_logger.NewLogger(loggerConfig)
	if err != nil {
		fmt.Println("logger cannot be initialized:", err)
		os.Exit(1)
	}
	logger.Debug("starting calculator-requests-reciever service")
	defer logger.Close()

	logger.Debug("initializing feature messages")
	messagesService := messages_service.NewMessagesService(nil)
	messagesHTTPHandler := messages_transport.NewMessagesHTTPHandler(&messagesService)

	logger.Debug("initializing HTTP server")
	apiVersion1Router := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersion1Router.RegisterRoutes(messagesHTTPHandler.Routes()...)

	serverConfig := core_http_server.NewConfigMust()
	server := core_http_server.NewHTTPServer(
		serverConfig,
		logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Recovery(),
		core_http_middleware.Trace(),
	)
	server.RegisterAPIRouters(apiVersion1Router)
	server.RegisterInfo()

	if err := server.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
