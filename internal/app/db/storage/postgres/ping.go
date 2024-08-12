package postgres

import (
	"context"
	"database/sql"
)

type PingStorage struct {
	db *sql.DB
}

func NewPingStorage(db *sql.DB) *PingStorage {
	return &PingStorage{db: db}
}

func (p *PingStorage) PingDB(ctx context.Context) error {
	return p.db.PingContext(ctx)
}
