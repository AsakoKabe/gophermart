package storage

import (
	"database/sql"

	"github.com/AsakoKabe/gophermart/internal/app/db/storage/postgres"
)

type Storages struct {
	PingStorage PingStorage
}

func NewPostgresStorages(db *sql.DB) (*Storages, error) {
	return &Storages{
		PingStorage: postgres.NewPingService(db),
	}, nil
}
