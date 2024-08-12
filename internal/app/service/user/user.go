package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/AsakoKabe/gophermart/internal/app/db/models"
	"github.com/AsakoKabe/gophermart/internal/app/db/storage"
	"github.com/AsakoKabe/gophermart/internal/app/db/storage/postgres"
	orderService "github.com/AsakoKabe/gophermart/internal/app/service/order"
)

type Service struct {
	userStorage       storage.UserStorage
	orderService      *orderService.Service
	withdrawalStorage storage.WithdrawalStorage
}

func NewService(
	userStorage storage.UserStorage,
	orderService *orderService.Service,
	withdrawalStorage storage.WithdrawalStorage,
) *Service {
	return &Service{
		userStorage:       userStorage,
		orderService:      orderService,
		withdrawalStorage: withdrawalStorage,
	}
}

var ErrLoginAlreadyExist = errors.New("user login already exist")

func (s *Service) Add(ctx context.Context, user *models.User) error {
	existedUser, err := s.userStorage.GetUserByLogin(ctx, user.Login)
	if err != nil && !errors.Is(err, postgres.ErrUserNotExist) {
		return err
	}

	if existedUser != nil {
		return ErrLoginAlreadyExist
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

func (s *Service) GetAccruals(ctx context.Context, userLogin string) (float64, error) {
	accruals, err := s.userStorage.GetAccruals(ctx, userLogin)
	if err != nil {
		slog.Error(
			"error to get accruals from storage", slog.String("err", err.Error()),
			slog.String("user login", userLogin),
		)
		return 0, err
	}

	return accruals, nil
}

func (s *Service) GetWithdrawal(ctx context.Context, userLogin string) (float64, error) {
	withdrawal, err := s.userStorage.GetWithdrawal(ctx, userLogin)
	if err != nil {
		slog.Error(
			"error to get withdrawal from storage", slog.String("err", err.Error()),
			slog.String("user login", userLogin),
		)
		return 0, err
	}

	return withdrawal, nil
}

func (s *Service) GetBalance(ctx context.Context, userLogin string) (float64, error) {
	balance, err := s.userStorage.GetBalance(ctx, userLogin)
	if err != nil {
		slog.Error(
			"error to get balance from storage", slog.String("err", err.Error()),
			slog.String("user login", userLogin),
		)
		return 0, err
	}

	return balance, nil
}

func (s *Service) GetAccrualsAndWithdrawal(ctx context.Context, userLogin string) (
	float64, float64, error,
) {
	accruals, withdrawal, err := s.userStorage.GetAccrualsAndWithdrawal(ctx, userLogin)
	if err != nil {
		slog.Error(
			"error to get accruals, withdrawal from storage", slog.String("err", err.Error()),
			slog.String("user login", userLogin),
		)
		return 0, 0, err
	}

	return accruals, withdrawal, nil
}
