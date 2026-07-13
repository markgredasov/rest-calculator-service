package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/markgredasov/rest-calculator-service/internal/app"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	app, err := app.InitializeApp(ctx)
	if err != nil {
		fmt.Printf("failed to initialize app: %w\n", err)
		os.Exit(1)
	}
	defer app.Close()

	if err := app.Server.Run(ctx); err != nil {
		app.Logger.Error("HTTP server run error", zap.Error(err))
	}
}
