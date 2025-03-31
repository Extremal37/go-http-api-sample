package psql

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	pgxv5 "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations
var fsMigrations embed.FS

func MigrateUp(ctx context.Context, postgresDsn string) error {
	sourceFs, err := iofs.New(fsMigrations, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create fs source driver: %w", err)
	}
	defer func() {
		if err = sourceFs.Close(); err != nil {
			fmt.Printf("failed to close file: %v", err)
		}
	}()

	pool, err := pgxpool.New(ctx, postgresDsn)
	if err != nil {
		return fmt.Errorf("unable to establish connection to database: %v", err)
	}
	defer func() {
		pool.Close()
	}()

	db := stdlib.OpenDB(*pool.Config().ConnConfig)
	conn, err := pgxv5.WithInstance(db, &pgxv5.Config{})
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceFs, "pgx", conn)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to UP migration: %w", err)
	}

	return nil
}
