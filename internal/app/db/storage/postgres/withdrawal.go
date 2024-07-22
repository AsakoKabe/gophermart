package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/AsakoKabe/gophermart/internal/app/db/models"
)

type WithdrawalStorage struct {
	db *sql.DB
}

func NewWithdrawalStorage(db *sql.DB) *WithdrawalStorage {
	return &WithdrawalStorage{db: db}
}

const insertWithdrawal = "insert into withdrawals (num_order, sum, user_id) VALUES ($1, $2, $3)"
const selectSumByUser = "select sum(sum) from withdrawals where user_id = $1"
const selectAllByUserID = "select num_order, sum, trim('\"' from to_json(processed_at)::text) from withdrawals where user_id = $1 order by processed_at"

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
	rows.Next()

	var sum sql.NullFloat64
	if err = rows.Scan(&sum); err != nil {
		slog.Error("error parse sum withdrawal from db", slog.String("err", err.Error()))
		return 0, err
	}
	if err = rows.Err(); err != nil {
		slog.Error("error to scan")
		return 0, err
	}
	if sum.Valid {
		return sum.Float64, nil
	}

	return 0, nil
}

func (s *WithdrawalStorage) GetAll(ctx context.Context, userID string) (*[]models.Withdrawal, error) {
	rows, err := s.db.QueryContext(
		ctx,
		selectAllByUserID,
		userID,
	)
	if err != nil {
		slog.Error("error select withdrawals by userID", slog.String("err", err.Error()))
		return nil, err
	}
	defer rows.Close()
	var withdrawals []models.Withdrawal

	for rows.Next() {
		withdrawal, errParse := s.parseWithdrawal(rows)
		if errParse != nil {
			continue
		}
		withdrawals = append(withdrawals, *withdrawal)
	}

	return &withdrawals, nil
}

func (s *WithdrawalStorage) parseWithdrawal(rows *sql.Rows) (*models.Withdrawal, error) {
	var withdrawal models.Withdrawal
	if err := rows.Scan(&withdrawal.OrderNum, &withdrawal.Sum, &withdrawal.ProcessedAt); err != nil {
		slog.Error("error parse withdrawal from db", slog.String("err", err.Error()))
		return nil, err
	}

	return &withdrawal, nil
}
