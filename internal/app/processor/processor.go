package processor

import (
	"github.com/Extremal37/go-http-api-sample/internal/app/storage"
	"go.uber.org/zap"
)

type Processor struct {
	log     *zap.SugaredLogger
	storage storage.Storage
}

func NewProcessor(storage storage.Storage, log *zap.SugaredLogger) *Processor {
	return &Processor{
		storage: storage,
		log:     log,
	}
}
