package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

func NewConfig() (config Config, err error) {
	config = Config{}
	if err = godotenv.Load(); err != nil {
		return
	}

	if err = env.Parse(&config); err != nil {
		return
	}

	return
}
