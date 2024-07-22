package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AsakoKabe/gophermart/internal/app/db/models"
	"log/slog"
)

type WithdrawalStorage struct {
	db *sql.DB
}

func NewWithdrawalStorage(db *sql.DB) *WithdrawalStorage {
	return &WithdrawalStorage{db: db}
}

const insertWithdrawal = "insert into withdrawals (num_order, sum, user_id) VALUES ($1, $2, $3)"
const selectSumByUser = "select sum(sum) from withdrawals where user_id = $1"

func (s *WithdrawalStorage) Add(ctx context.Context, withdrawal *models.Withdrawal) error {
	_, err := s.db.ExecContext(ctx, insertWithdrawal, withdrawal.OrderNum, withdrawal.Sum, withdrawal.UserID)
	if err != nil {
		return fmt.Errorf("unable to insert new withdrawal: %w", err)
	}

	return nil
}

func (s *WithdrawalStorage) GetSum(ctx context.Context, userID string) (float64, error) {
	rows, err := s.db.QueryContext(
		ctx,
		selectSumByUser,
		userID,
	)
	if err != nil {
		slog.Error("error select sum withdrawal by userID", slog.String("err", err.Error()))
		return 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, err
	}

	var sum sql.NullFloat64
	if err = rows.Scan(&sum); err != nil {
		slog.Error("error parse sum withdrawal from db", slog.String("err", err.Error()))
		return 0, err
	}
	if sum.Valid {
		return sum.Float64, nil
	}

	return 0, nil
}
