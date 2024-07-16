package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AuthSecret           string `env:"AUTH_SECRET"`
	Addr                 string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func LoadConfig() (*Config, error) {
	cfg := new(Config)

	parseFlag(cfg)

	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg, nil
}
