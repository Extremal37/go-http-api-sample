package handlers

import (
	"context"
	"github.com/Extremal37/go-http-api-sample/internal/app/models"
	"github.com/Extremal37/go-http-api-sample/internal/app/processor"
	"go.uber.org/zap"
)

type Processor interface {
	AddContact(ctx context.Context, contact models.Contact) error
	GetContacts(ctx context.Context) ([]models.Contact, error)
}

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
