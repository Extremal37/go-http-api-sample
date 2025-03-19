package processor

import "github.com/Extremal37/go-http-api-sample/internal/app/models"

// AddContact Saving new contact to the storage
func (p *Processor) AddContact(contact models.Contact) {
	p.storage.AddContact(contact)
}

// GetContacts Return all contacts from the storage
func (p *Processor) GetContacts() []models.Contact {
	return p.storage.GetContacts()
}
