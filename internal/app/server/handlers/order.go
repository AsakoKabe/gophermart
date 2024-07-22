package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/jwtauth/v5"

	"github.com/AsakoKabe/gophermart/internal/app/service"
	"github.com/AsakoKabe/gophermart/internal/app/service/order"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) Add(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	orderNum, err := readBody(r.Body)
	if err != nil {
		slog.Error("error to read body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())
	userLogin, ok := claims[tokenKey].(string)
	if !ok {
		slog.Error("error to get user login")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.orderService.Add(r.Context(), orderNum, userLogin)
	if err != nil {
		switch {
		case errors.Is(err, order.ErrAlreadyAddedOtherUser):
			w.WriteHeader(http.StatusConflict)
		case errors.Is(err, order.ErrAlreadyAdded):
			w.WriteHeader(http.StatusOK)
		case errors.Is(err, order.ErrBadFormat):
			w.WriteHeader(http.StatusUnprocessableEntity)
		default:
			slog.Error("error to add order", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	userLogin, ok := claims[tokenKey].(string)
	if !ok {
		slog.Error("error to get user login")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ordersAccrual, err := h.orderService.GetOrdersWithAccrual(r.Context(), userLogin)
	if err != nil {
		slog.Error("error to get orders with accrual", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if len(*ordersAccrual) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = json.NewEncoder(w).Encode(ordersAccrual)
	if err != nil {
		slog.Error("error to create response get orders", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func readBody(reqBody io.ReadCloser) (string, error) {
	body, err := io.ReadAll(reqBody)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	return string(body), nil
}
