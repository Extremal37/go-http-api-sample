package psql

import (
	"context"
	"fmt"
	"github.com/Extremal37/go-http-api-sample/internal/app/models"
	"github.com/jackc/pgx/v5"
)

type ContactStorageDTO struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
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

func (s *Storage) AddContact(ctx context.Context, contact models.Contact) error {
	c, cancel := context.WithTimeout(ctx, rowsRetrieveTimeout)
	defer cancel()

	dto := contactToStorage(contact)

	query := `INSERT INTO contacts (first_name,last_name) VALUES ($1,$2)`
	_, err := s.conn.Exec(c, query, dto.FirstName, dto.LastName)
	if err != nil {
		return fmt.Errorf("failed to insert contact to DB: %w", err)
	}
	return nil
}

func (s *Storage) GetContacts(ctx context.Context) ([]models.Contact, error) {
	c, cancel := context.WithTimeout(ctx, rowsRetrieveTimeout)
	defer cancel()

	query := `SELECT first_name,last_name FROM contacts`

	dtos := make([]ContactStorageDTO, 0)
	contacts := make([]models.Contact, 0)
	rows, err := s.conn.Query(c, query)
	if err != nil {
		return contacts, fmt.Errorf("failed to get contacts from DB: %w", err)
	}
	dtos, err = pgx.CollectRows(rows, pgx.RowToStructByName[ContactStorageDTO])
	if err != nil {
		return contacts, fmt.Errorf("failed to scan rows to struct: %w", err)
	}

	for _, v := range dtos {
		contacts = append(contacts, storageToContact(v))
	}

	return contacts, nil
}
