package server

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AsakoKabe/gophermart/internal/app/service"
	"github.com/go-chi/jwtauth/v5"
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
	dbPool   *sql.DB
	services *service.Services

	httpServer *http.Server
	tokenAuth  *jwtauth.JWTAuth
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
		dbPool:    pool,
		services:  service.NewServices(storages),
		tokenAuth: jwtauth.New("HS256", []byte(cfg.AuthSecret), nil),
	}, nil
}

func (a *App) Run(cfg *config.Config) error {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	err := a.registerHTTPEndpoint(router)
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

func (a *App) registerHTTPEndpoint(router *chi.Mux) error {
	pingHandler := handlers.NewPingHandler(a.services.PingService)
	router.Get("/ping", pingHandler.HealthDB)

	userHandler := handlers.NewUserHandler(a.services.UserService, a.tokenAuth)
	//orderHandler := handlers.NewOrderHandler(a.storages.OrderStorage)
	router.Route("/api/user/", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(a.tokenAuth))
			r.Use(jwtauth.Authenticator(a.tokenAuth))
			r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
				_, claims, _ := jwtauth.FromContext(r.Context())
				fmt.Println(claims)
			})

			//r.Post("/orders", orderHandler.Add)
		})
	})

	return nil
}
