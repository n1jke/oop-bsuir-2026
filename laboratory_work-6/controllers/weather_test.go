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
				client.EXPECT().LocationCurrentTemperature(53.9, 27.5667).Return(21.5, nil)
			},
			lat:     53.9,
			lon:     27.5667,
			want:    weather.CurrentWeather{Temperature: 21.5},
			wantErr: false,
		},
		{
			name: "client error",
			prepare: func(client *mock.MockWeatherDataClient) {
				client.EXPECT().LocationCurrentTemperature(10.0, 20.0).Return(0.0, errors.New("OGOGO"))
			},
			lat:     10.0,
			lon:     20.0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := mock.NewMockWeatherDataClient(ctrl)
			tt.prepare(client)
			c := controllers.NewCurrentWeatherController(client)

			got, gotErr := c.GetWeather(tt.lat, tt.lon)

			if tt.wantErr {
				require.Error(t, gotErr)
				return
			}

			require.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
