package storage

import (
	"context"

	"github.com/AsakoKabe/gophermart/internal/app/db/models"
)

type UserStorage interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	GetAccruals(ctx context.Context, userLogin string) (float64, error)
	GetWithdrawal(ctx context.Context, userLogin string) (float64, error)
	GetBalance(ctx context.Context, userLogin string) (float64, error)
	GetAccrualsAndWithdrawal(ctx context.Context, userLogin string) (
		float64, float64, error,
	)
}
