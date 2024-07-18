package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AsakoKabe/gophermart/internal/app/db/models"
	"log/slog"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{db: db}
}

const insertUser = "insert into users (login, password) values ($1, $2)"
const selectUser = "select * from users where login = $1"

func (u *UserStorage) CreateUser(ctx context.Context, user *models.User) error {
	_, err := u.db.ExecContext(ctx, insertUser, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("unable to insert new user: %w", err)
	}

	return nil
}

func (u *UserStorage) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	rows, err := u.db.QueryContext(
		ctx,
		selectUser,
		login,
	)
	if err != nil {
		slog.Error("error select user by login", slog.String("err", err.Error()))
		return nil, err
	}
	defer rows.Close()

	user, err := u.parseUser(rows)
	if err != nil {
		slog.Error("error to parse user", slog.String("err", err.Error()))
		return nil, err
	}

	return user, nil
}

func (u *UserStorage) parseUser(rows *sql.Rows) (*models.User, error) {
	if !rows.Next() {
		return nil, nil
	}
	var user models.User
	if err := rows.Scan(&user.ID, &user.Login, &user.Password); err != nil {
		slog.Error("error parse user from db", slog.String("err", err.Error()))
		return nil, err
	}

	return &user, nil
}
