package utils

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	OpenWeatherKey   string `env:"OPENWEATHER_API_KEY,required,notEmpty"`
	OpenWeatherURL   string `env:"OPENWEATHER_BASE_URL,required,notEmpty"`
	GoogleWeatherKey string `env:"GOOGLEWEATHER_API_KEY,required,notEmpty"`
	GoogleWeatherURL string `env:"GOOGLEWEATHER_BASE_URL,required,notEmpty"`
}

func LoadConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &AppConfig{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
