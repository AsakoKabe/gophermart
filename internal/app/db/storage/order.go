package storage

import "github.com/AsakoKabe/gophermart/internal/app/db/models"

type OrderStorage interface {
	Add(order *models.Order) error
}
