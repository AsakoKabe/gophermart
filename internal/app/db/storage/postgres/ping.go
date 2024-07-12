package postgres

import (
	"context"
	"database/sql"
)

type PingService struct {
	db *sql.DB
}

func NewPingService(db *sql.DB) *PingService {
	return &PingService{db: db}
}

func (p *PingService) PingDB(ctx context.Context) error {
	return p.db.PingContext(ctx)
}
