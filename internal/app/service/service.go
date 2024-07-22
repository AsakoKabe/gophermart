package service

import (
	"github.com/AsakoKabe/gophermart/config"
	"github.com/AsakoKabe/gophermart/internal/app/db/storage"
	"github.com/AsakoKabe/gophermart/internal/app/service/order"
	"github.com/AsakoKabe/gophermart/internal/app/service/ping"
	"github.com/AsakoKabe/gophermart/internal/app/service/user"
	"github.com/AsakoKabe/gophermart/internal/app/service/withdrawal"
)

type Services struct {
	UserService       UserService
	PingService       PingService
	OrderService      OrderService
	WithdrawalService WithdrawalService
}

func NewServices(storages *storage.Storages, cfg *config.Config) *Services {
	orderService := order.NewService(
		storages.OrderStorage,
		storages.UserStorage,
		cfg.AccrualSystemAddress,
	)
	return &Services{
		UserService:       user.NewService(storages.UserStorage, orderService, storages.WithdrawalStorage),
		PingService:       ping.NewService(storages.PingStorage),
		OrderService:      orderService,
		WithdrawalService: withdrawal.NewService(storages.WithdrawalStorage, storages.UserStorage),
	}
}
