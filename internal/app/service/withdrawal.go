package service

import (
	"context"

	"github.com/AsakoKabe/gophermart/internal/app/db/models"
)

type WithdrawalService interface {
	Add(ctx context.Context, orderNum string, sum float64, userLogin string) error
	GetAll(ctx context.Context, userLogin string) ([]*models.Withdrawal, error)
}
