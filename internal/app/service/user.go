package service

import (
	"context"

	"github.com/AsakoKabe/gophermart/internal/app/db/models"
)

type UserService interface {
	Add(ctx context.Context, user *models.User) error
	IsValidUser(ctx context.Context, user *models.User) (bool, error)
	GetAccruals(ctx context.Context, userLogin string) (float64, error)
	GetWithdrawal(ctx context.Context, userLogin string) (float64, error)
	GetBalance(ctx context.Context, userLogin string) (float64, error)
	GetAccrualsAndWithdrawal(ctx context.Context, userLogin string) (float64, float64, error)
}
