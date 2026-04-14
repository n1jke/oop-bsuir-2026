package controllers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/models/weather"
)

var ErrUnsupportedCity = errors.New("unsupported city")

//go:generate mockgen -source=weather.go -destination=mock/mock.go -package=mock
type WeatherDataClient interface {
	LocationCurrentTemperature(lat, lon float64) (weather.CurrentWeather, error)
	LocationForecast(lat, lon float64) (forecast weather.Forecast, err error)

	CityCurrentTemperature(city string) (weather.CurrentWeather, error)
	CityForecast(city string) (forecast weather.Forecast, err error)
}

type CurrentWeatherController[T WeatherDataClient] struct {
	Client T
}

func NewCurrentWeatherController[T WeatherDataClient](client T) *CurrentWeatherController[T] {
	return &CurrentWeatherController[T]{
		Client: client,
	}
}

func (c *CurrentWeatherController[T]) GetWeatherCoordinates(lat, lon float64) (weather.CurrentWeather, error) {
	result, err := c.Client.LocationCurrentTemperature(lat, lon)
	if err != nil {
		return weather.CurrentWeather{}, fmt.Errorf("fail to get current weather for coord: %w", err)
	}

	return result, nil
}

func (c *CurrentWeatherController[T]) GetWeatherCity(city string) (weather.CurrentWeather, error) {
	if !isAllowedCity(city) {
		return weather.CurrentWeather{}, fmt.Errorf("%w: %s", ErrUnsupportedCity, city)
	}

	result, err := c.Client.CityCurrentTemperature(city)
	if err != nil {
		return weather.CurrentWeather{}, fmt.Errorf("fail to get current weather for city %s: %w", city, err)
	}

	return result, nil
}

func (c *CurrentWeatherController[T]) GetForecastCoordinates(lat, lon float64) (weather.Forecast, error) {
	forecast, err := c.Client.LocationForecast(lat, lon)
	if err != nil {
		return weather.Forecast{}, fmt.Errorf("fail to get forecast for coordinates: %w", err)
	}

	return forecast, nil
}

func (c *CurrentWeatherController[T]) GetForecastCity(city string) (weather.Forecast, error) {
	if !isAllowedCity(city) {
		return weather.Forecast{}, fmt.Errorf("unsupported city: %s", city)
	}

	forecast, err := c.Client.CityForecast(city)
	if err != nil {
		return weather.Forecast{}, fmt.Errorf("fail to get forecast for city %s: %w", city, err)
	}

	return forecast, nil
}

func isAllowedCity(city string) bool {
	city = strings.ToLower(city)
	switch city {
	case "минск", "лондон", "токио", "шанхай", "варшава":
		return true
	default:
		return false
	}
}
