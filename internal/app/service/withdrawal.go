package service

import (
	"context"
)

type WithdrawalService interface {
	Add(ctx context.Context, orderNum string, sum float64, userLogin string) error
}
