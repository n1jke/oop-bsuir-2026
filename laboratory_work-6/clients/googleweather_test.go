package clients_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/clients"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/clients/mock"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/models/weather"
)

func TestGoogleWeatherClient_LocationCurrentTemperature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpClient := mock.NewMockWebClient(ctrl)
	client, err := clients.NewGoogleWeatherClient("gkey", "https://weather.googleapis.com/v1", httpClient)
	require.NoError(t, err)

	tests := []struct {
		name    string
		prepare func(client *mock.MockWebClient)
		lat     float64
		lon     float64
		want    weather.CurrentWeather
		wantErr bool
	}{
		{
			name: "ok",
			prepare: func(client *mock.MockWebClient) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"temperature":{"degrees":8.5}}`)),
				}
				client.EXPECT().
					Get("https://weather.googleapis.com/v1/currentConditions:lookup?key=gkey&location.latitude=53.900000&location.longitude=27.566700").Return(resp, nil)
			},
			lat:     53.9,
			lon:     27.5667,
			want:    weather.CurrentWeather{Temperature: 8.5},
			wantErr: false,
		},
		{
			name: "client error",
			prepare: func(client *mock.MockWebClient) {
				client.EXPECT().
					Get("https://weather.googleapis.com/v1/currentConditions:lookup?key=gkey&location.latitude=1.000000&location.longitude=2.000000").Return(nil, errors.New("network"))
			},
			lat:     1,
			lon:     2,
			wantErr: true,
		},
		{
			name: "!= 200 status",
			prepare: func(client *mock.MockWebClient) {
				resp := &http.Response{
					StatusCode: http.StatusBadGateway,
					Body:       io.NopCloser(bytes.NewBufferString(`{"temperature":{"degrees":8.5}}`)),
				}
				client.EXPECT().
					Get("https://weather.googleapis.com/v1/currentConditions:lookup?key=gkey&location.latitude=3.000000&location.longitude=4.000000").Return(resp, nil)
			},
			lat:     3,
			lon:     4,
			wantErr: true,
		},
		{
			name: "invalid json",
			prepare: func(client *mock.MockWebClient) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"temperature":{"degrees":"bad"}}`)),
				}
				client.EXPECT().
					Get("https://weather.googleapis.com/v1/currentConditions:lookup?key=gkey&location.latitude=5.000000&location.longitude=6.000000").Return(resp, nil)
			},
			lat:     5,
			lon:     6,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(httpClient)

			got, err := client.LocationCurrentTemperature(tt.lat, tt.lon)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGoogleWeatherClient_LocationForecast(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpClient := mock.NewMockWebClient(ctrl)
	client, err := clients.NewGoogleWeatherClient("gkey", "https://weather.googleapis.com/v1", httpClient)
	require.NoError(t, err)

	tests := []struct {
		name    string
		prepare func(client *mock.MockWebClient)
		lat     float64
		lon     float64
		want    weather.Forecast
		wantErr bool
	}{
		{
			name: "ok",
			prepare: func(client *mock.MockWebClient) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(
						`{"forecastDays":[` +
							`{"interval":{"startTime":"2025-02-10T15:00:00Z"},"maxTemperature":{"degrees":13.3}},` +
							`{"interval":{"startTime":"2025-02-11T15:00:00Z"},"maxTemperature":{"degrees":12.1}},` +
							`{"interval":{"startTime":"2025-02-12T15:00:00Z"},"maxTemperature":{"degrees":11.9}}` +
							`]}`)),
				}
				client.EXPECT().
					Get("https://weather.googleapis.com/v1/forecast/days:lookup?key=gkey&location.latitude=53.900000&location.longitude=27.566700&days=3").Return(resp, nil)
			},
			lat: 53.9,
			lon: 27.5667,
			want: weather.Forecast{Points: []weather.ForecastPoint{
				{Time: 1739199600, Temperature: 13.3},
				{Time: 1739286000, Temperature: 12.1},
				{Time: 1739372400, Temperature: 11.9},
			}},
			wantErr: false,
		},
		{
			name: "insufficient days",
			prepare: func(client *mock.MockWebClient) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(
						`{"forecastDays":[` +
							`{"interval":{"startTime":"2025-02-10T15:00:00Z"},"maxTemperature":{"degrees":13.3}},` +
							`{"interval":{"startTime":"2025-02-11T15:00:00Z"},"maxTemperature":{"degrees":12.1}}` +
							`]}`)),
				}
				client.EXPECT().
					Get("https://weather.googleapis.com/v1/forecast/days:lookup?key=gkey&location.latitude=1.000000&location.longitude=2.000000&days=3").Return(resp, nil)
			},
			lat:     1,
			lon:     2,
			wantErr: true,
		},
		{
			name: "invalid json",
			prepare: func(client *mock.MockWebClient) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"forecastDays":"bad"}`)),
				}
				client.EXPECT().
					Get("https://weather.googleapis.com/v1/forecast/days:lookup?key=gkey&location.latitude=3.000000&location.longitude=4.000000&days=3").Return(resp, nil)
			},
			lat:     3,
			lon:     4,
			wantErr: true,
		},
		{
			name: "!= 200 status",
			prepare: func(client *mock.MockWebClient) {
				resp := &http.Response{
					StatusCode: http.StatusBadGateway,
					Body:       io.NopCloser(bytes.NewBufferString(`{"forecastDays":[]}`)),
				}
				client.EXPECT().
					Get("https://weather.googleapis.com/v1/forecast/days:lookup?key=gkey&location.latitude=5.000000&location.longitude=6.000000&days=3").Return(resp, nil)
			},
			lat:     5,
			lon:     6,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(httpClient)

			got, err := client.LocationForecast(tt.lat, tt.lon)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
