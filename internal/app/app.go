package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Extremal37/go-http-api-sample/internal/app/processor"
	"github.com/Extremal37/go-http-api-sample/internal/cfg"
	"github.com/gorilla/mux"
	"net/http"

	"go.uber.org/zap"
)

const appName = "HTTP API Sample Server by Dmitry Tumalanov"

type Server struct {
	log    *zap.SugaredLogger
	config *cfg.Configuration

	httpServer *http.Server
	processor  *processor.Processor
	storage    processor.Storage
}

func NewServer(cfg *cfg.Configuration, processor *processor.Processor, storage processor.Storage, logger *zap.SugaredLogger) *Server {
	return &Server{
		log:       logger,
		processor: processor,
		storage:   storage,
		config:    cfg,
	}
}

func (s *Server) Serve(routes *mux.Router) {
	s.log.Infof("%s starting", appName)
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
