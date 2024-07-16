package postgres

import (
	"context"
	"database/sql"
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

const countLogin = "select count(*) from users where login = $1"
const insertUser = "insert into users (login, password) values ($1, $2)"
const selectUser = "select * from users where login = $1"

func (u *UserStorage) CreateUser(ctx context.Context, user *models.User) error {
	loginExist, err := u.isLoginExist(ctx, countLogin, user.Login)
	if err != nil {
		return err
	}

	if loginExist == true {
		return models.LoginAlreadyExist
	}

	_, err = u.db.ExecContext(ctx, insertUser, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("unable to insert new user: %w", err)
	}

	return nil
}

func (u *UserStorage) isLoginExist(ctx context.Context, query string, args ...any) (bool, error) {
	rows, err := u.db.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		slog.Error("error select count user by login", slog.String("err", err.Error()))
		return false, err
	}
	defer rows.Close()

	rows.Next()
	var count int
	if err = rows.Scan(&count); err != nil {
		slog.Error("error parse count user by login from db", slog.String("err", err.Error()))
		return false, err
	}

	if err = rows.Err(); err != nil {
		return false, err
	}

	return count != 0, nil
}

func (u *UserStorage) IsUserValid(ctx context.Context, user *models.User) (bool, error) {
	rows, err := u.db.QueryContext(
		ctx,
		selectUser,
		user.Login,
	)
	if err != nil {
		slog.Error("error select user by login", slog.String("err", err.Error()))
		return false, err
	}
	defer rows.Close()

	existedUser, err := u.parseUser(rows)
	if err != nil {
		slog.Error("error to parse user", slog.String("err", err.Error()))
		return false, err
	}
	if isEqual(existedUser, user) {
		return true, nil
	}
	return false, nil
}

func isEqual(existedUser *models.User, user *models.User) bool {
	return existedUser.Login == user.Login && existedUser.Password == user.Password
}

func (u *UserStorage) parseUser(rows *sql.Rows) (*models.User, error) {
	rows.Next()
	var user models.User
	if err := rows.Scan(&user.ID, &user.Login, &user.Password); err != nil {
		slog.Error("error parse user from db", slog.String("err", err.Error()))
		return nil, err
	}

	return &user, nil
}
