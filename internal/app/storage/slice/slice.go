package slice

import (
	"github.com/Extremal37/go-http-api-sample/internal/app/models"
	"go.uber.org/zap"
)

type Storage struct {
	slice []ContactStorageDTO
	log   *zap.SugaredLogger
}

type ContactStorageDTO struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

func NewStorage(log *zap.SugaredLogger) *Storage {
	slice := make([]ContactStorageDTO, 0)
	return &Storage{
		slice: slice,
		log:   log.With("storage", "slice"),
	}
}

func contactToStorage(contact models.Contact) ContactStorageDTO {
	return ContactStorageDTO{
		FirstName: contact.FirstName,
		LastName:  contact.LastName,
	}

}

func storageToContact(storage ContactStorageDTO) models.Contact {
	return models.Contact{
		FirstName: storage.FirstName,
		LastName:  storage.LastName,
	}
}

func (s *Storage) AddContact(contact models.Contact) {
	s.slice = append(s.slice, contactToStorage(contact))
	s.log.Infof("Contact added: %v", contact)
}

func (s *Storage) GetContacts() []models.Contact {
	var contacts []models.Contact

	for _, v := range s.slice {
		contacts = append(contacts, storageToContact(v))
	}
	return contacts
}
