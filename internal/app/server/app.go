package server

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/AsakoKabe/gophermart/config"
	"github.com/AsakoKabe/gophermart/internal/app/db/connection"
	"github.com/AsakoKabe/gophermart/internal/app/db/storage"
	"github.com/AsakoKabe/gophermart/internal/app/server/handlers"
)

const httpTimeOut = 10 * time.Second
const MaxHeaderBytes = 1 << 20
const ctxTimeout = 5 * time.Second

type App struct {
	httpServer *http.Server
	dbPool     *sql.DB
	storages   *storage.Storages
}

func NewApp(cfg *config.Config) (*App, error) {
	if cfg.DatabaseURI == "" {
		return nil, ErrConnectToDB
	}
	pool, err := connection.NewDBPool(cfg.DatabaseURI)
	if err != nil {
		slog.Error("error to create db pool", slog.String("err", err.Error()))
		return nil, ErrCreateDBPoll
	}

	storages, err := storage.NewPostgresStorages(pool)
	if err != nil {
		slog.Error("error to create service", slog.String("err", err.Error()))
		return nil, ErrCreateStorages
	}

	return &App{
		dbPool:   pool,
		storages: storages,
	}, nil
}

func (a *App) Run(cfg *config.Config) error {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	err := handlers.RegisterHTTPEndpoint(router, a.storages, cfg)
	if err != nil {
		return ErrRegisterEndpoints
	}

	a.httpServer = &http.Server{
		Addr:           cfg.Addr,
		Handler:        router,
		ReadTimeout:    httpTimeOut,
		WriteTimeout:   httpTimeOut,
		MaxHeaderBytes: MaxHeaderBytes,
	}

	go func() {
		err = a.httpServer.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func (a *App) CloseDBPool() {
	if a.dbPool == nil {
		return
	}
	err := a.dbPool.Close()
	if err != nil {
		return
	}
}