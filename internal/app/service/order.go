package service

import (
	"context"
)

type OrderService interface {
	Add(ctx context.Context, numOrder int, userLogin string) error
}
