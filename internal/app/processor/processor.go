package processor

import (
	"github.com/Extremal37/go-http-api-sample/internal/app/models"
	"go.uber.org/zap"
)

type Storage interface {
	AddContact(contact models.Contact)
	GetContacts() []models.Contact
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
