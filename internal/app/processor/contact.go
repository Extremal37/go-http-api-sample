package processor

import (
	"context"
	"github.com/Extremal37/go-http-api-sample/internal/app/models"
)

// AddContact Saving new contact to the storage
func (p *Processor) AddContact(ctx context.Context, contact models.Contact) error {
	return p.storage.AddContact(ctx, contact)
}

// GetContacts Return all contacts from the storage
func (p *Processor) GetContacts(ctx context.Context) ([]models.Contact, error) {
	return p.storage.GetContacts(ctx)
}
