package service

import (
	"context"

	"github.com/AsakoKabe/gophermart/internal/app/db/models"
)

type OrderService interface {
	Add(ctx context.Context, numOrder string, userLogin string) error
	GetOrders(ctx context.Context, userLogin string) (*[]models.Order, error)
	AddAccrualToOrders(ctx context.Context, orders *[]models.Order) (*[]models.OrderWithAccrual, error)
}
