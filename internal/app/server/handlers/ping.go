package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AsakoKabe/gophermart/internal/app/service"
)

type PingHandler struct {
	pingService service.PingService
}

func NewPingHandler(pingService service.PingService) *PingHandler {
	return &PingHandler{pingService: pingService}
}

func (h *PingHandler) HealthDB(w http.ResponseWriter, r *http.Request) {
	err := h.pingService.PingDB(r.Context())
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		message, _ := json.Marshal(map[string]string{"name": err.Error()})
		http.Error(w, string(message), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ok"))
}
