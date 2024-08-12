package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/AsakoKabe/gophermart/internal/app/db/models"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{db: db}
}

const insertUser = "insert into users (login, password) values ($1, $2)"
const selectUser = "select * from users where login = $1"
const selectAccrualsByLogin = "select accruals from users where login = $1"
const selectWithdrawalByLogin = "select withdrawal from users where login = $1"
const selectBalanceByLogin = "select (accruals-withdrawal) from users where login = $1"
const selectAccrualsAndWithdrawal = "select accruals, withdrawal from users where login = $1"

var ErrUserNotExist = errors.New("user not exist")

func (u *UserStorage) CreateUser(ctx context.Context, user *models.User) error {
	_, err := u.db.ExecContext(ctx, insertUser, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("unable to insert new user: %w", err)
	}

	return nil
}

func (u *UserStorage) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	row := u.db.QueryRowContext(
		ctx,
		selectUser,
		login,
	)
	user, err := u.parseUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotExist
	} else if err != nil {
		slog.Error("error to parse user", slog.String("err", err.Error()))
		return nil, err
	}

	return user, nil
}

func (u *UserStorage) parseUser(rows *sql.Row) (*models.User, error) {
	var user models.User
	if err := rows.Scan(
		&user.ID, &user.Login, &user.Password, &user.Accruals, &user.Withdrawal,
	); err != nil {
		slog.Error("error parse user from db", slog.String("err", err.Error()))
		return nil, err
	}

	return &user, nil
}

func (u *UserStorage) GetAccruals(ctx context.Context, userLogin string) (float64, error) {
	var accruals float64
	err := u.db.QueryRowContext(
		ctx,
		selectAccrualsByLogin,
		userLogin,
	).Scan(&accruals)

	if err != nil {
		slog.Error("error to select accruals", slog.String("err", err.Error()))
		return 0, err
	}
	return accruals, nil
}

func (u *UserStorage) GetWithdrawal(ctx context.Context, userLogin string) (float64, error) {
	var withdrawal float64
	err := u.db.QueryRowContext(
		ctx,
		selectWithdrawalByLogin,
		userLogin,
	).Scan(&withdrawal)

	if err != nil {
		slog.Error("error to select withdrawal", slog.String("err", err.Error()))
		return 0, err
	}
	return withdrawal, nil
}

func (u *UserStorage) GetBalance(ctx context.Context, userLogin string) (float64, error) {
	var balance float64
	err := u.db.QueryRowContext(
		ctx,
		selectBalanceByLogin,
		userLogin,
	).Scan(&balance)
	if err != nil {
		slog.Error("error to select balance", slog.String("err", err.Error()))
		return 0, err
	}
	return balance, nil
}

func (u *UserStorage) GetAccrualsAndWithdrawal(ctx context.Context, userLogin string) (
	float64, float64, error,
) {
	var accruals, withdrawal float64
	err := u.db.QueryRowContext(
		ctx,
		selectAccrualsAndWithdrawal,
		userLogin,
	).Scan(&accruals, &withdrawal)
	if err != nil {
		slog.Error("error to select accruals, withdrawal", slog.String("err", err.Error()))
		return 0, 0, err
	}
	return accruals, withdrawal, nil
}
