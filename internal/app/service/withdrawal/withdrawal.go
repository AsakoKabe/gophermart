package withdrawal

import (
	"context"
	"log/slog"

	"github.com/AsakoKabe/gophermart/internal/app/db/models"
	"github.com/AsakoKabe/gophermart/internal/app/db/storage"
)

type Service struct {
	withdrawalStorage storage.WithdrawalStorage
	userStorage       storage.UserStorage
}

func NewService(withdrawalStorage storage.WithdrawalStorage, userStorage storage.UserStorage) *Service {
	return &Service{withdrawalStorage: withdrawalStorage, userStorage: userStorage}
}

func (s *Service) Add(ctx context.Context, orderNum string, sum float64, userLogin string) error {
	user, err := s.userStorage.GetUserByLogin(ctx, userLogin)
	if err != nil {
		slog.Error("error to get user for add withdrawal",
			slog.String("userLogin", userLogin),
			slog.String("err", err.Error()),
		)
		return err
	}

	err = s.withdrawalStorage.Add(ctx, &models.Withdrawal{
		OrderNum: orderNum,
		Sum:      sum,
		UserID:   user.ID,
	})
	if err != nil {
		slog.Error("error to add withdrawal", slog.String("err", err.Error()))
		return err
	}

	return nil
}
