package app

import (
	"context"
	"errors"
	"github.com/Extremal37/go-http-api-sample/internal/app/processor"
	"github.com/Extremal37/go-http-api-sample/internal/cfg"
	"github.com/gorilla/mux"
	"net/http"

	"go.uber.org/zap"
)

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
	s.log.Infof("Starting HTTP listener on address %s", s.config.App.Address)
	s.httpServer = &http.Server{
		Addr:    s.config.App.Address,
		Handler: routes,
	}

	if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		s.log.Fatalf("Failed to start http server: %v", err)
	}

}
func (s *Server) Shutdown() {
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(context.Background()); err != nil {
			s.log.Errorf("Failed to shutdown HTTP Server: %v", err)
		}
	}
}
