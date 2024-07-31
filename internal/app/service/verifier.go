package service

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/AsakoKabe/gophermart/internal/app/db/models"
	"github.com/AsakoKabe/gophermart/internal/app/db/storage"
	"github.com/go-resty/resty/v2"
)

type AccrualResponse struct {
	Order   string             `json:"order"`
	Status  models.OrderStatus `json:"status"`
	Accrual float64            `json:"accrual"`
}

var ErrTooManyRequestsAccrual = errors.New("too many request in accrual service")
var ErrNotRegisterOrderAccrual = errors.New("order isnt registered in accrual")
var ErrUnexpectedStatusCode = errors.New("unexpected status code in accrual")

const accrualURI = "/api/orders/"
const duration = 100 * time.Millisecond

type Verifier struct {
	httpClient   *resty.Client
	orderStorage storage.OrderStorage
	userStorage  storage.UserStorage
	accrualURI   string
	ticker       *time.Ticker
	done         chan bool
}

func NewVerifier(
	orderStorage storage.OrderStorage, userStorage storage.UserStorage, accrualURI string,
) *Verifier {
	return &Verifier{
		httpClient: resty.New(), orderStorage: orderStorage, userStorage: userStorage,
		accrualURI: accrualURI,
		ticker:     time.NewTicker(duration),
		done:       make(chan bool),
	}
}

func (v *Verifier) Start() {
	slog.Info("start verifier")
	go v.update()
}

func (v *Verifier) Stop() {
	slog.Info("stop verifier")
	v.ticker.Stop()
	v.done <- true
}

func (v *Verifier) update() {
	for {
		select {
		case <-v.done:
			return
		case <-v.ticker.C:
			unprocessedOrders := v.getUnprocessedOrders()
			for _, order := range unprocessedOrders {
				v.updateOrderStatus(order)
				v.updateAccruals(order)
			}
		}
	}
}

func (v *Verifier) updateOrderStatus(order *models.Order) {
	err := v.orderStorage.UpdateOrderStatus(
		context.Background(), order.Status, order.Num,
	)
	if err != nil {
		slog.Error(
			"error to update order status",
			slog.String("err", err.Error()),
			slog.String("orderNum", order.Num),
			slog.Any("newStatus", order.Status),
		)
	}
}

func (v *Verifier) updateAccruals(order *models.Order) {
	if order.Status == models.PROCESSED {
		err := v.userStorage.UpdateAccruals(
			context.Background(), order.UserID, order.Accrual,
		)
		if err != nil {
			slog.Error(
				"error to update user accruals",
				slog.String("err", err.Error()),
				slog.String("user id", order.UserID),
				slog.Any("accrual", order.Accrual),
			)
			return
		}
		err = v.orderStorage.UpdateOrderAccrual(
			context.Background(), order.ID, order.Accrual,
		)
		if err != nil {
			slog.Error(
				"error to update order accrual",
				slog.String("err", err.Error()),
				slog.String("order id", order.ID),
				slog.Any("accrual", order.Accrual),
			)
			return
		}
	}
}

func (v *Verifier) getUnprocessedOrders() []*models.Order {
	orders, err := v.orderStorage.GetOrdersWithStatuses(
		context.Background(),
		[]models.OrderStatus{models.NEW, models.PROCESSING},
	)
	if err != nil {
		slog.Error("error to get orders with statuses", slog.String("err", err.Error()))
		return []*models.Order{}
	}

	var unprocessedOrders []*models.Order

	for _, order := range orders {
		var errResponse error
		accrualResponse, errResponse := v.sendAccrualRequest(order.Num)
		if errResponse != nil {
			switch {
			case errors.Is(errResponse, ErrTooManyRequestsAccrual):
				slog.Error(
					errResponse.Error(),
					slog.String("orderNum", order.Num),
				)
				continue
			default:
				slog.Error(
					"error to get accrualResponse for order",
					slog.String("order num", order.Num),
					slog.String("err", errResponse.Error()),
				)
				continue
			}
		}

		unprocessedOrders = append(
			unprocessedOrders, &models.Order{
				ID:         order.ID,
				Num:        order.Num,
				Status:     accrualResponse.Status,
				Accrual:    accrualResponse.Accrual,
				UploadedAt: order.UploadedAt,
				UserID:     order.UserID,
			},
		)
	}

	return unprocessedOrders
}

func (v *Verifier) sendAccrualRequest(orderNum string) (*AccrualResponse, error) {
	var accrualResponse AccrualResponse
	resp, err := v.httpClient.R().
		SetResult(&accrualResponse).
		ForceContentType("application/json").
		Get(v.accrualURI + accrualURI + orderNum)

	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusTooManyRequests:
		return nil, ErrTooManyRequestsAccrual
	case http.StatusNoContent:
		accrualResponse.Status = models.NEW
		return &accrualResponse, nil
	case http.StatusInternalServerError:
		return nil, err
	case http.StatusOK:
		return &accrualResponse, nil
	}

	return nil, ErrUnexpectedStatusCode
}
