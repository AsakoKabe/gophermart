package storage

import (
	"database/sql"

	"github.com/AsakoKabe/gophermart/internal/app/db/storage/postgres"
)

type Storages struct {
	PingStorage  PingStorage
	UserStorage  UserStorage
	OrderStorage OrderStorage
}

func NewPostgresStorages(db *sql.DB) (*Storages, error) {
	return &Storages{
		PingStorage:  postgres.NewPingStorage(db),
		UserStorage:  postgres.NewUserStorage(db),
		OrderStorage: postgres.NewOrderStorage(db),
	}, nil
}
