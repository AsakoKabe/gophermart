package handlers

import (
	"encoding/json"
	"errors"
	"github.com/AsakoKabe/gophermart/internal/app/service/order"
	"github.com/go-chi/jwtauth/v5"
	"log/slog"
	"net/http"

	"github.com/AsakoKabe/gophermart/internal/app/service"
)

type WithdrawalHandler struct {
	withdrawalService service.WithdrawalService
	orderService      service.OrderService
	userService       service.UserService
}

func NewWithdrawalHandler(
	withdrawalService service.WithdrawalService,
	orderService service.OrderService,
	userService service.UserService,
) *WithdrawalHandler {
	return &WithdrawalHandler{
		withdrawalService: withdrawalService,
		orderService:      orderService,
		userService:       userService,
	}
}

func (h *WithdrawalHandler) Add(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var withdrawal withdrawalRequest
	err := json.NewDecoder(r.Body).Decode(&withdrawal)
	if err != nil {
		slog.Error("error to read body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())
	userLogin, ok := claims[tokenKey].(string)
	if !ok {
		slog.Error("error to get user login", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balance, err := h.userService.GetBalance(r.Context(), userLogin)
	if err != nil {
		slog.Error(
			"error to get balance for add withdrawal",
			slog.String("userLogin", userLogin),
			slog.String("err", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if balance < withdrawal.Sum {
		slog.Error("user doesnt have enough balance", slog.String("userLogin", userLogin))
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

	err = h.orderService.Add(r.Context(), withdrawal.Order, userLogin)
	if err != nil {
		if errors.Is(err, order.ErrBadFormat) {
			slog.Error("bad format order num")
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		slog.Error("error to add order")
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = h.withdrawalService.Add(r.Context(), withdrawal.Order, withdrawal.Sum, userLogin)
	if err != nil {
		slog.Error(
			"error to add withdrawal",
			slog.Any("withdrawal", withdrawal),
			slog.String("err", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
