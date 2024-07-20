package storage

import (
	"context"

	"github.com/AsakoKabe/gophermart/internal/app/db/models"
)

type OrderStorage interface {
	Add(ctx context.Context, order *models.Order) error
	GetOrderByNum(ctx context.Context, num string) (*models.Order, error)
	GetOrdersByUserIDSortedByUpdatedAt(ctx context.Context, userID string) (*[]models.Order, error)
}
