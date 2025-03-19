package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Extremal37/go-http-api-sample/api"
	"github.com/Extremal37/go-http-api-sample/internal/app/handlers"
	"github.com/Extremal37/go-http-api-sample/internal/app/processor"
	"github.com/Extremal37/go-http-api-sample/internal/app/storage"
	"github.com/Extremal37/go-http-api-sample/internal/app/storage/slice"
	"github.com/Extremal37/go-http-api-sample/internal/cfg"
	"net/http"

	"go.uber.org/zap"
)

const appName = "HTTP API Sample Server by Dmitry Tumalanov"

type Server struct {
	log    *zap.SugaredLogger
	config *cfg.Configuration

	httpServer *http.Server
	processor  *processor.Processor
	storage    storage.Storage
}

func NewServer(cfg *cfg.Configuration, logger *zap.SugaredLogger) *Server {
	return &Server{
		log:    logger,
		config: cfg,
	}
}

func (s *Server) Serve() {
	s.log.Infof("%s starting", appName)

	s.log.Debug("Connecting to storage")
	s.storage = slice.NewStorage(s.log)

	s.log.Debug("Spawning processor and handler")
	s.processor = processor.NewProcessor(s.storage, s.log)

	hdl := handlers.NewHandler(s.processor, s.log)

	routes := api.CreateRoutes(hdl, s.log)

	s.log.Infof("Starting HTTP listener on port %d", s.config.App.ListenPort)
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.App.ListenPort),
		Handler: routes,
	}

	if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		s.log.Fatalf("Failed to start http server: %v", err)
	}

}
func (s *Server) Shutdown() {
	if s.httpServer != nil {
		_ = s.httpServer.Shutdown(context.Background())
	}
}
