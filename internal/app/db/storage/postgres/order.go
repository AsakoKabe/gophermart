package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/AsakoKabe/gophermart/internal/app/db/models"
)

type OrderStorage struct {
	db *sql.DB
}

func NewOrderStorage(db *sql.DB) *OrderStorage {
	return &OrderStorage{db: db}
}

const insertOrder = "insert into orders (num, user_id) values ($1, $2)"
const selectOrderByNum = "select id, num, user_id, trim('\"' from to_json(uploaded_at)::text) from orders where num = $1"
const selectOrdersByUserID = "select id, num, user_id, trim('\"' from to_json(uploaded_at)::text) from orders where user_id = $1 order by uploaded_at"

func (s *OrderStorage) Add(ctx context.Context, order *models.Order) error {
	_, err := s.db.ExecContext(ctx, insertOrder, order.Num, order.UserID)
	if err != nil {
		return fmt.Errorf("unable to insert new user: %w", err)
	}

	return nil
}

func (s *OrderStorage) GetOrderByNum(ctx context.Context, num string) (*models.Order, error) {
	rows, err := s.db.QueryContext(
		ctx,
		selectOrderByNum,
		num,
	)
	if err != nil {
		slog.Error("error select order by num", slog.String("err", err.Error()))
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}

	order, err := s.parseOrder(rows)
	if err != nil {
		slog.Error("error to parse order", slog.String("err", err.Error()))
		return nil, err
	}

	return order, nil
}

func (s *OrderStorage) parseOrder(rows *sql.Rows) (*models.Order, error) {
	var order models.Order
	if err := rows.Scan(&order.ID, &order.Num, &order.UserID, &order.UploadedAt); err != nil {
		slog.Error("error parse order from db", slog.String("err", err.Error()))
		return nil, err
	}

	return &order, nil
}

func (s *OrderStorage) GetOrdersByUserIDSortedByUpdatedAt(ctx context.Context, userID string) (*[]models.Order, error) {
	rows, err := s.db.QueryContext(
		ctx,
		selectOrdersByUserID,
		userID,
	)
	if err != nil {
		slog.Error("error select order by userID", slog.String("err", err.Error()))
		return nil, err
	}
	defer rows.Close()
	var orders []models.Order

	for rows.Next() {
		order, errParse := s.parseOrder(rows)
		if errParse != nil {
			continue
		}
		orders = append(orders, *order)
	}

	return &orders, nil
}
