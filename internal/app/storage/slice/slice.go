package slice

import (
	"github.com/Extremal37/go-http-api-sample/internal/app/models"
	"go.uber.org/zap"
)

type Storage struct {
	slice *[]models.Contact
	log   *zap.SugaredLogger
}

func NewStorage(log *zap.SugaredLogger) *Storage {
	slice := make([]models.Contact, 0)
	return &Storage{
		slice: &slice,
		log:   log.With("storage", "slice"),
	}
}

func (s *Storage) AddContact(contact models.Contact) {
	*s.slice = append(*s.slice, contact)
	s.log.Infof("Contact added: %v", contact)
}

func (s *Storage) GetContacts() *[]models.Contact {
	return s.slice
}
