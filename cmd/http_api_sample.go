package main

import (
	"context"
	"github.com/Extremal37/go-http-api-sample/api"
	"github.com/Extremal37/go-http-api-sample/internal/app"
	"github.com/Extremal37/go-http-api-sample/internal/app/handlers"
	"github.com/Extremal37/go-http-api-sample/internal/app/processor"
	"github.com/Extremal37/go-http-api-sample/internal/app/storage/slice"
	"github.com/Extremal37/go-http-api-sample/internal/cfg"
	"github.com/Extremal37/go-http-api-sample/internal/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	// Load config
	config, err := cfg.LoadAndStoreConfig()
	if err != nil {
		log.Fatalf("Cannot load config - %v", err)
	}

	// Init logger
	logger := log.NewLogger(config.App.Logging)

	// Creating server with loaded config
	logger.Debug("Connecting to storage")
	storage := slice.NewStorage(logger)

	logger.Debug("Spawning processor and handler")
	proc := processor.NewProcessor(storage, logger)
	hdl := handlers.NewHandler(proc, logger)
	routes := api.CreateRoutes(hdl, logger)

	server := app.NewServer(config, proc, storage, logger)

	// Launching server

	go server.Serve(routes)

	// Ждём сигнала завершения приложения
	select {
	case <-ctx.Done():
		stop()
		server.Shutdown()
		logger.Info("The service has been terminated successfully")
		os.Exit(0)
	}

}
