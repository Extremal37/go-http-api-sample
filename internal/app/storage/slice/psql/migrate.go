package psql

import (
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	pgxv5 "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations
var fsMigrations embed.FS

func (s *Storage) migrateUp() error {
	sourceFs, err := iofs.New(fsMigrations, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create fs source driver: %w", err)
	}
	db := stdlib.OpenDB(*s.conn.Config().ConnConfig)
	conn, err := pgxv5.WithInstance(db, &pgxv5.Config{})
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	defer func() {
		if err = sourceFs.Close(); err != nil {
			s.log.Errorf("failed to close file: %v", err)
		}
	}()
	m, err := migrate.NewWithInstance("iofs", sourceFs, "pgx", conn)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
