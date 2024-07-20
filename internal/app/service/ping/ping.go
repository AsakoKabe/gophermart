package ping

import (
	"context"

	"github.com/AsakoKabe/gophermart/internal/app/db/storage"
)

type Service struct {
	pingStorage storage.PingStorage
}

func NewService(pingStorage storage.PingStorage) *Service {
	return &Service{pingStorage: pingStorage}
}

func (s *Service) PingDB(ctx context.Context) error {
	return s.pingStorage.PingDB(ctx)
}
