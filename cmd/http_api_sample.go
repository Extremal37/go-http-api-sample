package main

import (
	"context"
	"github.com/Extremal37/go-http-api-sample/internal/app"
	"github.com/Extremal37/go-http-api-sample/internal/cfg"
	"github.com/Extremal37/go-http-api-sample/internal/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load config
	config, err := cfg.LoadAndStoreConfig()
	if err != nil {
		log.Fatalf("Cannot load config - %v", err)
	}

	// Init logger
	logger := log.NewLogger(config.App.Logging)

	// Creating server with loaded config
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	server := app.NewServer(config, logger)

	// Launching server

	go server.Serve()

	// Ждём сигнала завершения приложения
	select {
	case <-ctx.Done():
		stop()
		server.Shutdown()
		logger.Info("The service has been terminated successfully")
		os.Exit(0)
	}

}
