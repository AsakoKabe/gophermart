package storage

import "context"

type User struct {
}

type UserStorage interface {
	CreateUser(ctx context.Context, user User) (string, error)
}
