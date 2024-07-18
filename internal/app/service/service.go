package service

import (
	"github.com/AsakoKabe/gophermart/internal/app/db/storage"
	"github.com/AsakoKabe/gophermart/internal/app/service/ping"
	"github.com/AsakoKabe/gophermart/internal/app/service/user"
)

type Services struct {
	UserService UserService
	PingService PingService
}

func NewServices(storages *storage.Storages) *Services {
	return &Services{
		UserService: user.NewService(storages.UserStorage),
		PingService: ping.NewService(storages.PingStorage),
	}
}
