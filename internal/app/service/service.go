package service

import (
	"github.com/AsakoKabe/gophermart/config"
	"github.com/AsakoKabe/gophermart/internal/app/db/storage"
	"github.com/AsakoKabe/gophermart/internal/app/service/order"
	"github.com/AsakoKabe/gophermart/internal/app/service/ping"
	"github.com/AsakoKabe/gophermart/internal/app/service/user"
)

type Services struct {
	UserService  UserService
	PingService  PingService
	OrderService OrderService
}

func NewServices(storages *storage.Storages, cfg *config.Config) *Services {
	return &Services{
		UserService: user.NewService(storages.UserStorage),
		PingService: ping.NewService(storages.PingStorage),
		OrderService: order.NewService(
			storages.OrderStorage,
			storages.UserStorage,
			cfg.AccrualSystemAddress,
		),
	}
}
