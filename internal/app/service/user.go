package service

import (
	"context"

	"github.com/AsakoKabe/gophermart/internal/app/db/models"
)

type UserService interface {
	Add(ctx context.Context, user *models.User) error
	IsValidUser(ctx context.Context, user *models.User) (bool, error)
}
