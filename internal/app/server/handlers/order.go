package handlers

import (
	"errors"
	"github.com/AsakoKabe/gophermart/internal/app/service"
	"github.com/AsakoKabe/gophermart/internal/app/service/order"
	"github.com/go-chi/jwtauth/v5"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) Add(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	numOrder, err := readBody(r.Body)
	if err != nil {
		slog.Error("error to read body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	num, err := strconv.Atoi(numOrder)
	if err != nil {
		slog.Error(
			"error with type num order",
			slog.String("num order", numOrder),
			slog.String("err", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())
	userLogin := claims[tokenKey].(string)

	err = h.orderService.Add(r.Context(), num, userLogin)
	if errors.Is(err, order.AlreadyAddedOtherUser) {
		w.WriteHeader(http.StatusConflict)
	} else if errors.Is(err, order.AlreadyAdded) {
		w.WriteHeader(http.StatusOK)
	} else if errors.Is(err, order.BadFormat) {
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else if err != nil {
		slog.Error("error to add order", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}

}

func readBody(reqBody io.ReadCloser) (string, error) {
	body, err := io.ReadAll(reqBody)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	return string(body), nil
}
