package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AsakoKabe/gophermart/internal/app/db/storage"
)

type PingHandler struct {
	pingStorage storage.PingStorage
}

func NewPingHandler(pingStorage storage.PingStorage) *PingHandler {
	return &PingHandler{pingStorage: pingStorage}
}

func (h *PingHandler) HealthDB(w http.ResponseWriter, r *http.Request) {
	err := h.pingStorage.PingDB(r.Context())
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		message, _ := json.Marshal(map[string]string{"name": err.Error()})
		http.Error(w, string(message), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ok"))
}
