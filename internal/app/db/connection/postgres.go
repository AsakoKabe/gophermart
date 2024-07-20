package connection

import (
	"database/sql"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"log"
)

//go:embed migrations/*.sql
var migrationsFolder embed.FS

func NewDBPool(dsn string) (*sql.DB, error) {
	pool, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func RunMigrations(dsn string) error {
	fs, err := iofs.New(migrationsFolder, "migrations")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", fs, dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Steps(2)
	if err != nil {
		return err
	}

	return nil
}
