package postgres

import (
	"database/sql"
	"github.com/AsakoKabe/gophermart/internal/app/db/models"
)

type OrderStorage struct {
	db *sql.DB
}

func NewOrderStorage(db *sql.DB) *OrderStorage {
	return &OrderStorage{db: db}
}

func (o *OrderStorage) Add(order *models.Order) error {
	//TODO implement me
	panic("implement me")
}
