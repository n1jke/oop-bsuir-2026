package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct{}

func (c *AppConfig) Validate() error {
	return nil
}

func LoadEnv() error {
	return godotenv.Load()
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
