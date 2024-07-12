package storage

import "context"

type PingStorage interface {
	PingDB(ctx context.Context) error
}
