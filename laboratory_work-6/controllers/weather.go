package controllers

import (
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-6/clients"
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-6/models/weather"

	"github.com/shopspring/decimal"
)

type CurrentWeatherController[T clients.WeatherDataClient] struct {
	Client T
}

func NewCurrentWeatherController[T clients.WeatherDataClient](client T) *CurrentWeatherController[T] {
	return &CurrentWeatherController[T]{
		Client: client,
	}
}

func (c *CurrentWeatherController[T]) GetCurrentWeather(lat decimal.Decimal, lon decimal.Decimal) (weather.CurrentWeather, error) {
	temperature, err := c.Client.LocationCurrentTemperature(lat, lon)
	if err != nil {
		return weather.CurrentWeather{}, err
	}

	return weather.CurrentWeather{
		Temperature: temperature,
	}, nil
}
