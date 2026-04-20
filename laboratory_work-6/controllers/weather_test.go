package controllers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/controllers"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/controllers/mock"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/models/weather"
)

func TestCurrentWeatherController_GetWeather(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockWeatherDataClient(ctrl)
	c := controllers.NewCurrentWeatherController(client)

	tests := []struct {
		name    string
		prepare func(client *mock.MockWeatherDataClient)
		lat     float64
		lon     float64
		want    weather.CurrentWeather
		wantErr bool
	}{
		{
			name: "ok",
			prepare: func(client *mock.MockWeatherDataClient) {
				client.EXPECT().LocationCurrentTemperature(53.9, 27.5667).Return(
					weather.CurrentWeather{Temperature: 21.5},
					nil,
				)
			},
			lat:     53.9,
			lon:     27.5667,
			want:    weather.CurrentWeather{Temperature: 21.5},
			wantErr: false,
		},
		{
			name: "client error",
			prepare: func(client *mock.MockWeatherDataClient) {
				client.EXPECT().LocationCurrentTemperature(10.0, 20.0).Return(
					weather.CurrentWeather{},
					errors.New("OGOGO"),
				)
			},
			lat:     10.0,
			lon:     20.0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(client)

			got, gotErr := c.GetWeatherCoordinates(tt.lat, tt.lon)

			if tt.wantErr {
				require.Error(t, gotErr)
				return
			}

			require.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCurrentWeatherController_GetWeatherCity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockWeatherDataClient(ctrl)
	c := controllers.NewCurrentWeatherController(client)

	tests := []struct {
		name    string
		prepare func(client *mock.MockWeatherDataClient)
		city    string
		want    weather.CurrentWeather
		wantErr bool
	}{
		{
			name: "ok",
			prepare: func(client *mock.MockWeatherDataClient) {
				client.EXPECT().CityCurrentTemperature("Минск").Return(
					weather.CurrentWeather{Temperature: 21.5},
					nil,
				)
			},
			city:    "Минск",
			want:    weather.CurrentWeather{Temperature: 21.5},
			wantErr: false,
		},
		{
			name: "invalid city",
			prepare: func(_ *mock.MockWeatherDataClient) {
			},
			city:    "москва",
			wantErr: true,
		},
		{
			name: "client error",
			prepare: func(client *mock.MockWeatherDataClient) {
				client.EXPECT().CityCurrentTemperature("Лондон").Return(
					weather.CurrentWeather{},
					errors.New("some err"),
				)
			},
			city:    "Лондон",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(client)

			got, gotErr := c.GetWeatherCity(tt.city)

			if tt.wantErr {
				require.Error(t, gotErr)
				return
			}

			require.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCurrentWeatherController_GetForecastCoordinates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockWeatherDataClient(ctrl)
	c := controllers.NewCurrentWeatherController(client)

	tests := []struct {
		name    string
		prepare func(client *mock.MockWeatherDataClient)
		lat     float64
		lon     float64
		want    weather.Forecast
		wantErr bool
	}{
		{
			name: "ok",
			prepare: func(client *mock.MockWeatherDataClient) {
				client.EXPECT().LocationForecast(53.9, 27.5667).Return(
					weather.Forecast{Points: []weather.ForecastPoint{
						{Time: 1776157200, Temperature: 9.3},
					}},
					nil,
				)
			},
			lat: 53.9,
			lon: 27.5667,
			want: weather.Forecast{Points: []weather.ForecastPoint{
				{Time: 1776157200, Temperature: 9.3},
			}},
			wantErr: false,
		},
		{
			name: "client error",
			prepare: func(client *mock.MockWeatherDataClient) {
				client.EXPECT().LocationForecast(10.0, 20.0).Return(weather.Forecast{}, errors.New("some err"))
			},
			lat:     10.0,
			lon:     20.0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(client)

			got, gotErr := c.GetForecastCoordinates(tt.lat, tt.lon)

			if tt.wantErr {
				require.Error(t, gotErr)
				return
			}

			require.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCurrentWeatherController_GetForecastCity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockWeatherDataClient(ctrl)
	c := controllers.NewCurrentWeatherController(client)

	tests := []struct {
		name    string
		prepare func(client *mock.MockWeatherDataClient)
		city    string
		want    weather.Forecast
		wantErr bool
	}{
		{
			name: "ok",
			prepare: func(client *mock.MockWeatherDataClient) {
				client.EXPECT().CityForecast("Токио").Return(
					weather.Forecast{Points: []weather.ForecastPoint{
						{Time: 1713024000, Temperature: 4.8},
					}},
					nil,
				)
			},
			city: "Токио",
			want: weather.Forecast{Points: []weather.ForecastPoint{
				{Time: 1713024000, Temperature: 4.8},
			}},
			wantErr: false,
		},
		{
			name: "invalid city",
			prepare: func(_ *mock.MockWeatherDataClient) {
			},
			city:    "Берлин",
			wantErr: true,
		},
		{
			name: "client error",
			prepare: func(client *mock.MockWeatherDataClient) {
				client.EXPECT().CityForecast("Варшава").Return(weather.Forecast{}, errors.New("some err"))
			},
			city:    "Варшава",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(client)

			got, gotErr := c.GetForecastCity(tt.city)

			if tt.wantErr {
				require.Error(t, gotErr)
				return
			}

			require.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
