package utils

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/clients"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/controllers"
)

func ParseProvideFlag(cfg *AppConfig) (controllers.WeatherDataClient, error) {
	provider := flag.String("provide", "google", "Weather provider, supported: google, openweather")

	flag.Parse()

	switch *provider {
	case "google":
		return clients.NewGoogleWeatherClient(cfg.GoogleWeatherKey, cfg.GoogleWeatherURL, http.DefaultClient)
	case "openweather":
		return clients.NewOpenWeatherClient(cfg.OpenWeatherKey, cfg.OpenWeatherURL, http.DefaultClient)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", *provider)
	}
}
