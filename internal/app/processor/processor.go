package processor

import (
	"context"
	"github.com/Extremal37/go-http-api-sample/internal/app/models"
	"go.uber.org/zap"
)

type Storage interface {
	AddContact(ctx context.Context, contact models.Contact) error
	GetContacts(ctx context.Context) ([]models.Contact, error)
	Stop()
}

type Processor struct {
	log     *zap.SugaredLogger
	storage Storage
}

func NewProcessor(storage Storage, log *zap.SugaredLogger) *Processor {
	return &Processor{
		storage: storage,
		log:     log,
	}
}
