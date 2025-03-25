package psql

import (
	"context"
	"fmt"
	"github.com/Extremal37/go-http-api-sample/internal/cfg"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"net"
	"time"
)

const (
	driverName          = "postgres"
	connectionTimeout   = 5 * time.Second
	rowsRetrieveTimeout = 10 * time.Second
)

type Storage struct {
	conn *pgxpool.Pool
	log  *zap.SugaredLogger
}

func NewStorage(ctx context.Context, conf cfg.Postgres, log *zap.SugaredLogger) *Storage {
	endpoint := net.JoinHostPort(conf.Host, fmt.Sprintf("%d", conf.Port))
	postgresDsn := fmt.Sprintf("%s://%s:%s@%s/%s", driverName, conf.Username, conf.Password, endpoint, conf.Database)

	s := Storage{
		conn: nil,
		log:  log.With("postgres", endpoint),
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	db, err := pgxpool.New(ctxTimeout, postgresDsn)
	if err != nil {
		s.log.Fatalf("unable to establish connection to database: %s", err.Error())
		return nil
	}
	s.conn = db

	if err = s.migrateUp(); err != nil {
		s.conn.Close()
		s.log.Fatalf("unable to migrate up: %s", err.Error())
		return nil
	}

	s.log.Debug("Successfully connected")
	return &s
}

func (s *Storage) Stop() {
	s.conn.Close()
}
