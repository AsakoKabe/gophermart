package order

import (
	"context"
	"errors"
	"github.com/AsakoKabe/gophermart/internal/app/db/models"
	"github.com/AsakoKabe/gophermart/internal/app/db/storage"
	"log/slog"
)

type Service struct {
	orderStorage storage.OrderStorage
	userStorage  storage.UserStorage
}

func NewService(orderStorage storage.OrderStorage, userStorage storage.UserStorage) *Service {
	return &Service{orderStorage: orderStorage, userStorage: userStorage}
}

var AlreadyAddedOtherUser = errors.New("order already added other user")
var AlreadyAdded = errors.New("order already added")
var BadFormat = errors.New("bad format num order")

func (s *Service) Add(ctx context.Context, numOrder int, userLogin string) error {
	user, err := s.userStorage.GetUserByLogin(ctx, userLogin)
	if err != nil {
		slog.Error("error to get user for add order",
			slog.String("userLogin", userLogin),
			slog.String("err", err.Error()),
		)
		return err
	}

	existedOrder, err := s.orderStorage.GetOrderByNum(ctx, numOrder)
	if existedOrder != nil {
		if existedOrder.UserID != user.ID {
			return AlreadyAddedOtherUser
		}
		return AlreadyAdded
	}
	if err != nil {
		slog.Error("error to select order", slog.String("err", err.Error()))
		return err
	}

	err = s.orderStorage.Add(ctx, &models.Order{
		Num:    numOrder,
		UserID: user.ID,
	})
	if err != nil {
		slog.Error("error to add order", slog.String("err", err.Error()))
		return err
	}

	return nil
}
