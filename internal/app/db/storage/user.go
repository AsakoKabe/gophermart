package storage

import (
	"context"
	"github.com/AsakoKabe/gophermart/internal/app/db/models"
)

type UserStorage interface {
	CreateUser(ctx context.Context, user *models.User) error
	IsUserValid(ctx context.Context, user *models.User) (bool, error)
}
