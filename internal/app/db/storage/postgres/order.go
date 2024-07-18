package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AsakoKabe/gophermart/internal/app/db/models"
	"log/slog"
)

type OrderStorage struct {
	db *sql.DB
}

func NewOrderStorage(db *sql.DB) *OrderStorage {
	return &OrderStorage{db: db}
}

const insertOrder = "insert into orders (num, user_id) values ($1, $2)"
const selectOrder = "select * from orders where num = $1"

func (s *OrderStorage) Add(ctx context.Context, order *models.Order) error {
	_, err := s.db.ExecContext(ctx, insertOrder, order.Num, order.UserID)
	if err != nil {
		return fmt.Errorf("unable to insert new user: %w", err)
	}

	return nil
}

func (s *OrderStorage) GetOrderByNum(ctx context.Context, num int) (*models.Order, error) {
	rows, err := s.db.QueryContext(
		ctx,
		selectOrder,
		num,
	)
	if err != nil {
		slog.Error("error select order by num", slog.String("err", err.Error()))
		return nil, err
	}
	defer rows.Close()

	order, err := s.parseOrder(rows)
	if err != nil {
		slog.Error("error to parse order", slog.String("err", err.Error()))
		return nil, err
	}

	return order, nil
}

func (s *OrderStorage) parseOrder(rows *sql.Rows) (*models.Order, error) {
	if !rows.Next() {
		return nil, nil
	}
	var order models.Order
	if err := rows.Scan(&order.ID, &order.Num, &order.UserID); err != nil {
		slog.Error("error parse order from db", slog.String("err", err.Error()))
		return nil, err
	}

	return &order, nil
}
