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

func (s *Service) GetBalance(ctx context.Context, userLogin string) (float64, error) {
	ordersAccrual, err := s.orderService.GetOrdersWithAccrual(ctx, userLogin)
	if err != nil {
		slog.Error(
			"error to get orders with accrual",
			slog.String("userLogin", userLogin),
			slog.String("err", err.Error()),
		)
		return 0, err
	}

	var balance float64
	for _, order := range *ordersAccrual {
		balance += order.Accrual
	}

	return balance, nil
}

func (s *Service) GetSumWithdrawal(ctx context.Context, userLogin string) (float64, error) {
	user, err := s.userStorage.GetUserByLogin(ctx, userLogin)
	if err != nil {
		slog.Error("error to get user", slog.String("err", err.Error()))
		return 0, err
	}

	withdrawal, err := s.withdrawalStorage.GetSum(ctx, user.ID)
	if err != nil {
		slog.Error("error to get sum withdrawal")
		return 0, err
	}

	return withdrawal, nil
}
