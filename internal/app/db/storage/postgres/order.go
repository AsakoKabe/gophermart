package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/AsakoKabe/gophermart/internal/app/db/models"
)

type OrderStorage struct {
	db *sql.DB
}

func NewOrderStorage(db *sql.DB) *OrderStorage {
	return &OrderStorage{db: db}
}

type sqlRow interface {
	Scan(dest ...any) error
}

const insertOrder = "insert into orders (num, user_id) values ($1, $2)"
const selectOrderByNum = "select id, num, status, accrual, user_id, trim('\"' from to_json(uploaded_at)::text) from orders where num = $1"
const selectOrdersByUserID = "select id, num, status, accrual, user_id, trim('\"' from to_json(uploaded_at)::text) from orders where user_id = $1 order by uploaded_at"
const selectOrdersWithStatuses = "select id, num, status, accrual, user_id, trim('\"' from to_json(uploaded_at)::text) from orders where status in (%s)"
const updateOrderStatus = "update orders SET status = $1 where num = $2"
const updateOrderAccrual = "update orders set accrual=$1 where id=$2"

func (s *OrderStorage) Add(ctx context.Context, order *models.Order) error {
	_, err := s.db.ExecContext(ctx, insertOrder, order.Num, order.UserID)
	if err != nil {
		return fmt.Errorf("unable to insert new user: %w", err)
	}

	return nil
}

func (s *OrderStorage) GetOrderByNum(ctx context.Context, num string) (*models.Order, error) {
	row := s.db.QueryRowContext(
		ctx,
		selectOrderByNum,
		num,
	)

	order, err := parseOrder(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		slog.Error("error to parse order", slog.String("err", err.Error()))
		return nil, err
	}

	return order, nil
}

func parseOrder(row sqlRow) (*models.Order, error) {
	var order models.Order
	if err := row.Scan(
		&order.ID, &order.Num, &order.Status, &order.Accrual, &order.UserID, &order.UploadedAt,
	); err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *OrderStorage) GetOrdersByUserIDSortedByUpdatedAt(
	ctx context.Context, userID string,
) ([]*models.Order, error) {
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
	var orders []*models.Order

	for rows.Next() {
		order, errParse := parseOrder(rows)
		if errParse != nil {
			continue
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *OrderStorage) GetOrdersWithStatuses(
	ctx context.Context, statuses []models.OrderStatus,
) ([]*models.Order, error) {
	params := make([]string, 0, len(statuses))
	var vals []any

	for i := range statuses {
		params = append(params, fmt.Sprintf("$%d", i+1))
		vals = append(vals, statuses[i])
	}
	query := fmt.Sprintf(
		selectOrdersWithStatuses,
		strings.Join(params, ", "),
	)

	rows, err := s.db.QueryContext(
		ctx,
		query,
		vals...,
	)
	if err != nil {
		slog.Error(
			"error select orders with status",
			slog.String("err", err.Error()),
			slog.Any("statuses", statuses),
		)
		return nil, err
	}
	defer rows.Close()
	var orders []*models.Order

	for rows.Next() {
		order, errParse := parseOrder(rows)
		if errParse != nil {
			continue
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *OrderStorage) UpdateOrderStatus(
	ctx context.Context, newStatus models.OrderStatus, orderNum string,
) error {
	_, err := s.db.ExecContext(ctx, updateOrderStatus, newStatus, orderNum)
	if err != nil {
		return fmt.Errorf("unable to update order status: %w", err)
	}

	return nil
}

func (s *OrderStorage) UpdateOrderAccrual(
	ctx context.Context, orderID string, accrual float64,
) error {
	_, err := s.db.ExecContext(ctx, updateOrderAccrual, accrual, orderID)
	if err != nil {
		return fmt.Errorf("unable to update order accrual: %w", err)
	}

	return nil
}
