package controllers

import (
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/models/weather"
)

//go:generate mockgen -source=weather.go -destination=mock/mock.go -package=mock
type WeatherDataClient interface {
	LocationCurrentTemperature(lat, lon float64) (temperature float64, err error) // todo: return domain model instead of primitive
	LocationForecast(lat, lon float64) (forecast weather.Forecast, err error)

	CityCurrentTemperature(city string) (temperature float64, err error)
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
	temperature, err := c.Client.LocationCurrentTemperature(lat, lon)
	if err != nil {
		return weather.CurrentWeather{}, err
	}

	return weather.CurrentWeather{
		Temperature: temperature,
	}, nil
}

func (c *CurrentWeatherController[T]) GetWeatherCity(city string) (weather.CurrentWeather, error) {
	panic("not implemented")
}

func (c *CurrentWeatherController[T]) GetForecastCoordinates(lat, lon float64) (weather.Forecast, error) {
	panic("not implemented")
}

func (c *CurrentWeatherController[T]) GetForecastCity(city string) (weather.Forecast, error) {
	panic("not implemented")
}
