package service

import "context"

type PingService interface {
	PingDB(ctx context.Context) error
}
