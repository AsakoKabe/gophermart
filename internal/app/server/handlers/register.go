package handlers

import (
	"github.com/AsakoKabe/gophermart/config"
	"github.com/AsakoKabe/gophermart/internal/app/db/storage"
	"github.com/go-chi/chi/v5"
)

func RegisterHTTPEndpoint(router *chi.Mux, storages *storage.Storages, cfg *config.Config) error {
	pingHandler := NewPingHandler(storages.PingStorage)
	router.Get("/ping", pingHandler.healthDB)

	return nil
}
