package main

import (
	"context"
	"fmt"
	"github.com/Extremal37/go-http-api-sample/api"
	"github.com/Extremal37/go-http-api-sample/internal/app"
	"github.com/Extremal37/go-http-api-sample/internal/app/handlers"
	"github.com/Extremal37/go-http-api-sample/internal/app/processor"
	"github.com/Extremal37/go-http-api-sample/internal/app/storage/psql"
	"github.com/Extremal37/go-http-api-sample/internal/app/storage/slice"
	"github.com/Extremal37/go-http-api-sample/internal/cfg"
	"github.com/Extremal37/go-http-api-sample/internal/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"net"
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
		log.Fatalf("Cannot load config - %v", err.Error())
	}

	// Init logger
	logger, err := log.NewLogger(config.App.Logging)
	if err != nil {
		log.Fatalf("Cannot init logger: %v", err.Error())
	}

	// Creating server with loaded config
	logger.Infof("%s starting", appName)
	logger.Debugf("Connecting to storage %s", config.App.Storage)

	var storage processor.Storage
	switch config.App.Storage {
	case storageSlice:
		storage = slice.NewStorage(logger)
	case storagePostgres:
		endpoint := net.JoinHostPort(config.Postgres.Host, fmt.Sprintf("%d", config.Postgres.Port))
		postgresDsn := fmt.Sprintf("%s://%s:%s@%s/%s", psql.DriverName, config.Postgres.Username, config.Postgres.Password, endpoint, config.Postgres.Database)
		ctxTimeout, cancel := context.WithTimeout(ctx, psql.ConnectionTimeout)
		defer cancel()

		db, err := pgxpool.New(ctxTimeout, postgresDsn)
		if err != nil {
			log.Fatalf("Unable to establish connection to database: %v", err.Error())
		}
		defer func() {
			db.Close()
		}()

		if err = psql.MigrateUp(db); err != nil {
			db.Close()
			log.Fatalf("unable to migrate up: %s", err.Error())
		}

		logger.Info("Successfully connected")
		storage = psql.NewStorage(db, logger.With(storagePostgres, endpoint))
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
