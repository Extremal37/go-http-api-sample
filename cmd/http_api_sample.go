package main

import (
	"context"
	"github.com/Extremal37/go-http-api-sample/api"
	"github.com/Extremal37/go-http-api-sample/internal/app"
	"github.com/Extremal37/go-http-api-sample/internal/app/handlers"
	"github.com/Extremal37/go-http-api-sample/internal/app/processor"
	"github.com/Extremal37/go-http-api-sample/internal/app/storage/psql"
	"github.com/Extremal37/go-http-api-sample/internal/app/storage/slice"
	"github.com/Extremal37/go-http-api-sample/internal/cfg"
	"github.com/Extremal37/go-http-api-sample/internal/log"
	"os"
	"os/signal"
	"syscall"
)

const (
	appName         = "HTTP API Sample Server by Dmitry Tumalanov"
	storageSlice    = "slice"
	storagePostgres = "postgres"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	// Load config
	config, err := cfg.LoadAndStoreConfig()
	if err != nil {
		log.Fatalf("Cannot load config - %v", err)
	}

	// Init logger
	logger, err := log.NewLogger(config.App.Logging)
	if err != nil {
		log.Fatalf("Cannot init logger: %v", err)
	}

	// Creating server with loaded config
	logger.Infof("%s starting", appName)
	logger.Debugf("Connecting to storage %s", config.App.Storage)

	var storage processor.Storage
	switch config.App.Storage {
	case storageSlice:
		storage = slice.NewStorage(logger)
	case storagePostgres:
		storage = psql.NewStorage(ctx, config.Postgres, logger)
	default:
		logger.Fatalf("Unknown storage %s. Supported storages are %v", config.App.Storage, []string{storagePostgres, storageSlice})
	}

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
		err = server.Shutdown()
		if err != nil {
			logger.Errorf("The service has been terminated with error:%v", err)
			os.Exit(1)
		}
		logger.Info("The service has been terminated successfully")
		os.Exit(0)
	}

}
