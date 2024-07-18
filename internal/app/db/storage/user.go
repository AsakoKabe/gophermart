package storage

import (
	"context"
	"github.com/AsakoKabe/gophermart/internal/app/db/models"
)

type UserStorage interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
}
