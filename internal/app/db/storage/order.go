package storage

import (
	"context"

	"github.com/AsakoKabe/gophermart/internal/app/db/models"
)

type OrderStorage interface {
	Add(ctx context.Context, order *models.Order) error
	GetOrderByNum(ctx context.Context, num string) (*models.Order, error)
	GetOrdersByUserIDSortedByUpdatedAt(ctx context.Context, userID string) ([]*models.Order, error)
	GetOrdersWithStatuses(ctx context.Context, statuses []models.OrderStatus) (
		[]*models.Order, error,
	)
	UpdateAccrualAndStatus(
		ctx context.Context, orderID string, accrual float64, orderStatus models.OrderStatus,
		userID string,
	) error
}
