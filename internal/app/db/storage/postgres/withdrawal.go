package postgres

import (
	"context"
	"database/sql"
	"errors"
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
const selectBalanceByUserID = "select accruals-users.withdrawal from users where id=$1 for update"
const updateUserWithdrawal = "update users set withdrawal=withdrawal+$1 where id=$2"

var ErrBalanceNotEnough = errors.New("user balance doesnt enough")

func (s *WithdrawalStorage) Add(ctx context.Context, withdrawal *models.Withdrawal) error {
	tx, err := s.db.Begin()
	if err != nil {
		slog.Error("error to create transaction", slog.String("err", err.Error()))
		return err
	}
	defer tx.Commit()
	var balance float64
	err = tx.QueryRow(
		selectBalanceByUserID, withdrawal.UserID,
	).Scan(&balance)
	if err != nil {
		slog.Error("error to get balance", slog.String("err", err.Error()))
		return err
	}

	if balance < withdrawal.Sum {
		return ErrBalanceNotEnough
	}

	_, err = tx.ExecContext(
		ctx, insertWithdrawal, withdrawal.OrderNum, withdrawal.Sum, withdrawal.UserID,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("unable to insert new withdrawal: %w", err)
	}

	_, err = tx.ExecContext(ctx, updateUserWithdrawal, withdrawal.Sum, withdrawal.UserID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("unable to update user withdrawal: %w", err)
	}

	return nil
}

func (s *WithdrawalStorage) GetSum(ctx context.Context, userID string) (float64, error) {
	row := s.db.QueryRowContext(
		ctx,
		selectSumByUser,
		userID,
	)
	var sum sql.NullFloat64
	if err := row.Scan(&sum); err != nil {
		slog.Error("error parse sum withdrawal from db", slog.String("err", err.Error()))
		return 0, err
	}
	if err := row.Err(); err != nil {
		slog.Error("error to scan")
		return 0, err
	}
	if sum.Valid {
		return sum.Float64, nil
	}

	return 0, nil
}

func (s *WithdrawalStorage) GetAll(ctx context.Context, userID string) (
	[]*models.Withdrawal, error,
) {
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
	var withdrawals []*models.Withdrawal

	for rows.Next() {
		withdrawal, errParse := s.parseWithdrawal(rows)
		if errParse != nil {
			continue
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	return withdrawals, nil
}

func (s *WithdrawalStorage) parseWithdrawal(rows *sql.Rows) (*models.Withdrawal, error) {
	var withdrawal models.Withdrawal
	if err := rows.Scan(
		&withdrawal.OrderNum, &withdrawal.Sum, &withdrawal.ProcessedAt,
	); err != nil {
		slog.Error("error parse withdrawal from db", slog.String("err", err.Error()))
		return nil, err
	}

	return &withdrawal, nil
}
