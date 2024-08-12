package storage

import (
	"context"

	"github.com/AsakoKabe/gophermart/internal/app/db/models"
)

type WithdrawalStorage interface {
	Add(ctx context.Context, withdrawal *models.Withdrawal) error
	GetSum(ctx context.Context, userID string) (float64, error)
	GetAll(ctx context.Context, userID string) ([]*models.Withdrawal, error)
}
