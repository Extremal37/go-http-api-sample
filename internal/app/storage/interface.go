package storage

import "github.com/Extremal37/go-http-api-sample/internal/app/models"

type Storage interface {
	AddContact(contact models.Contact)
	GetContacts() []models.Contact
}
