package controllers

import (
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/models/weather"
)

//go:generate mockgen -source=weather.go -destination=mock/mock.go -package=mock
type WeatherDataClient interface {
	LocationCurrentTemperature(lat, lon float64) (temperature float64, err error)
}

type CurrentWeatherController[T WeatherDataClient] struct {
	Client T
}

func NewCurrentWeatherController[T WeatherDataClient](client T) *CurrentWeatherController[T] {
	return &CurrentWeatherController[T]{
		Client: client,
	}
}

func (c *CurrentWeatherController[T]) GetWeather(lat, lon float64) (weather.CurrentWeather, error) {
	temperature, err := c.Client.LocationCurrentTemperature(lat, lon)
	if err != nil {
		return weather.CurrentWeather{}, err
	}

	return weather.CurrentWeather{
		Temperature: temperature,
	}, nil
}
