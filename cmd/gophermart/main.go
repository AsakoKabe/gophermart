package main

import (
	"log"

	"github.com/AsakoKabe/gophermart/config"
	"github.com/AsakoKabe/gophermart/internal/app/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	app, err := server.NewApp(cfg)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	defer app.CloseDBPool()

	if err = app.Run(cfg); err != nil {
		log.Panicf("%s", err.Error())
	}
}
