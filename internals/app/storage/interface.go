package storage

import "github.com/Extremal37/go-http-api-sample/internals/app/models"

type Storage interface {
	AddContact(contact models.Contact)
	GetContacts() *[]models.Contact
}
