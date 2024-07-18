package user

import (
	"context"
	"errors"
	"github.com/AsakoKabe/gophermart/internal/app/db/models"
	"github.com/AsakoKabe/gophermart/internal/app/db/storage"
	"log/slog"
)

type Service struct {
	userStorage storage.UserStorage
}

func NewService(userStorage storage.UserStorage) *Service {
	return &Service{userStorage: userStorage}
}

var LoginAlreadyExist = errors.New("user login already exist")

func (s *Service) Add(ctx context.Context, user *models.User) error {
	existedUser, err := s.userStorage.GetUserByLogin(ctx, user.Login)
	if err != nil {
		return err
	}

	if existedUser != nil {
		return LoginAlreadyExist
	}

	err = s.userStorage.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) IsValidUser(ctx context.Context, user *models.User) (bool, error) {
	existedUser, err := s.userStorage.GetUserByLogin(ctx, user.Login)
	if err != nil {
		slog.Error("error to get user", slog.String("err", err.Error()))
		return false, err
	}
	return isEqual(existedUser, user), nil
}

func isEqual(existedUser *models.User, user *models.User) bool {
	return existedUser.Login == user.Login && existedUser.Password == user.Password
}