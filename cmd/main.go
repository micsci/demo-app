package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/demo-marketplace/internal"
	"github.com/demo-marketplace/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("loading config", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := internal.NewServer(logger, cfg)
	if err := server.Run(ctx); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
