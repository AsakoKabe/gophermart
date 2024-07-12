package connection

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewDBPool(dsn string) (*sql.DB, error) {
	pool, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
