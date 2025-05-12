package psql

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

const (
	DriverName          = "postgres"
	ConnectionTimeout   = 5 * time.Second
	rowsRetrieveTimeout = 10 * time.Second
)

type Storage struct {
	conn *pgxpool.Pool
	log  *zap.SugaredLogger
}

func NewStorage(conn *pgxpool.Pool, log *zap.SugaredLogger) *Storage {
	return &Storage{
		conn: conn,
		log:  log,
	}

}
