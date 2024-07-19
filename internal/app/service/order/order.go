package order

import (
	"context"
	"errors"
	"github.com/AsakoKabe/gophermart/internal/app/db/models"
	"github.com/AsakoKabe/gophermart/internal/app/db/storage"
	"github.com/go-resty/resty/v2"
	"log/slog"
	"net/http"
	"strconv"
)

type Service struct {
	orderStorage storage.OrderStorage
	userStorage  storage.UserStorage
	httpClient   *resty.Client
	accrualURI   string
}

func NewService(
	orderStorage storage.OrderStorage,
	userStorage storage.UserStorage,
	accrualURI string,
) *Service {
	return &Service{
		orderStorage: orderStorage,
		userStorage:  userStorage,
		httpClient:   resty.New(),
		accrualURI:   accrualURI,
	}
}

var AlreadyAddedOtherUser = errors.New("order already added other user")
var AlreadyAdded = errors.New("order already added")
var BadFormat = errors.New("bad format num order")
var TooManyRequestsAccrual = errors.New("too many request in accrual service")
var NotRegisterOrderAccrual = errors.New("order isnt registered in accrual")

const accrualURI = "/api/orders/"

func (s *Service) Add(ctx context.Context, numOrder int, userLogin string) error {
	user, err := s.userStorage.GetUserByLogin(ctx, userLogin)
	if err != nil {
		slog.Error("error to get user for add order",
			slog.String("userLogin", userLogin),
			slog.String("err", err.Error()),
		)
		return err
	}

	existedOrder, err := s.orderStorage.GetOrderByNum(ctx, numOrder)
	if existedOrder != nil {
		if existedOrder.UserID != user.ID {
			return AlreadyAddedOtherUser
		}
		return AlreadyAdded
	}
	if err != nil {
		slog.Error("error to select order", slog.String("err", err.Error()))
		return err
	}

	err = s.orderStorage.Add(ctx, &models.Order{
		Num:    numOrder,
		UserID: user.ID,
	})
	if err != nil {
		slog.Error("error to add order", slog.String("err", err.Error()))
		return err
	}

	return nil
}

func (s *Service) GetOrders(ctx context.Context, userLogin string) (*[]models.Order, error) {
	user, err := s.userStorage.GetUserByLogin(ctx, userLogin)
	if err != nil {
		slog.Error("error to get user for get orders",
			slog.String("userLogin", userLogin),
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	orders, err := s.orderStorage.GetOrdersByUserIDSortedByUpdatedAt(ctx, user.ID)
	if err != nil {
		slog.Error("error to get order from storage")
		return nil, err
	}
	return orders, nil
}

type AccrualResponse struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}

func (s *Service) AddAccrualToOrders(
	_ context.Context,
	orders *[]models.Order,
) (*[]models.OrderWithAccrual, error) {
	var ordersAccrual []models.OrderWithAccrual

	for _, order := range *orders {
		accrualResponse, err := s.sendAccrualResponse(order.Num)
		if errors.Is(err, TooManyRequestsAccrual) {
			slog.Info(
				err.Error(),
				slog.Int("orderNum", order.Num),
			)
			continue
		} else if errors.Is(err, NotRegisterOrderAccrual) {
			slog.Info(
				err.Error(),
				slog.Int("orderNum", order.Num),
			)
			continue
		} else if err != nil {
			slog.Error(
				"error to get accrualResponse for order",
				slog.Int("order num", order.Num),
				slog.String("err", err.Error()),
			)
			continue
		}

		ordersAccrual = append(ordersAccrual, models.OrderWithAccrual{
			Number:     strconv.Itoa(order.Num),
			Status:     accrualResponse.Status,
			Accrual:    accrualResponse.Accrual,
			UploadedAt: order.UploadedAt,
		})
	}

	return &ordersAccrual, nil
}

func (s *Service) sendAccrualResponse(orderNum int) (*AccrualResponse, error) {
	var accrualResponse AccrualResponse
	resp, err := s.httpClient.R().
		SetResult(&accrualResponse).
		ForceContentType("application/json").
		Get(s.accrualURI + accrualURI + strconv.Itoa(orderNum))

	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusUnprocessableEntity:
		return nil, TooManyRequestsAccrual
	case http.StatusNoContent:
		return nil, NotRegisterOrderAccrual
	case http.StatusInternalServerError:
		return nil, err
	case http.StatusOK:
		return &accrualResponse, nil
	}

	return nil, nil

}
