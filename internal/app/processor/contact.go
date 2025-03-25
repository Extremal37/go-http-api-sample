package processor

import (
	"context"
	"fmt"
	"github.com/Extremal37/go-http-api-sample/internal/app/models"
)

// AddContact Saving new contact to the storage
func (p *Processor) AddContact(ctx context.Context, contact models.Contact) error {
	err := p.storage.AddContact(ctx, contact)
	if err != nil {
		return fmt.Errorf("failed to add contact to storage: %w", err)
	}
	return nil
}

// GetContacts Return all contacts from the storage
func (p *Processor) GetContacts(ctx context.Context) ([]models.Contact, error) {
	contacts, err := p.storage.GetContacts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts from storage: %w", err)
	}
	return contacts, nil
}
