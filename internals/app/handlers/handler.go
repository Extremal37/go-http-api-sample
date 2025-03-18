package handlers

import (
	"github.com/Extremal37/go-http-api-sample/internals/app/processor"
	"go.uber.org/zap"
)

type Handler struct {
	p   *processor.Processor
	log *zap.SugaredLogger
}

func NewHandler(p *processor.Processor, log *zap.SugaredLogger) *Handler {
	return &Handler{
		p:   p,
		log: log,
	}
}
