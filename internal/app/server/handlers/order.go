package handlers

import (
	"net/http"

	"github.com/AsakoKabe/gophermart/internal/app/db/storage"
)

type OrderHandler struct {
	orderStorage storage.OrderStorage
}

func NewOrderHandler(orderStorage storage.OrderStorage) *OrderHandler {
	return &OrderHandler{orderStorage: orderStorage}
}

func (h *OrderHandler) Add(w http.ResponseWriter, r *http.Request) {

}
