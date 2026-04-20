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

func TestNewOpenWeatherClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockWebClient(ctrl)

	tests := []struct {
		name    string
		apiKey  string
		baseURL string
		client  clients.WebClient
		wantErr bool
	}{
		{
			name:    "without api key",
			apiKey:  "",
			baseURL: "api/v1",
			client:  client,
			wantErr: true,
		},
		{
			name:    "without base url",
			apiKey:  "afosij913",
			baseURL: "",
			client:  client,
			wantErr: true,
		},
		{
			name:    "without client",
			apiKey:  "amdf928",
			baseURL: "api/v2",
			client:  nil,
			wantErr: true,
		},
		{
			name:    "valid",
			apiKey:  "ma9FG8F",
			baseURL: "api/v1",
			client:  client,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := clients.NewOpenWeatherClient(tt.apiKey, tt.baseURL, tt.client)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOpenWeatherClient_LocationCurrentTemperature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpClient := mock.NewMockWebClient(ctrl)
	client, err := clients.NewOpenWeatherClient("8rfne8", "api/v1", httpClient)
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
					Body:       io.NopCloser(bytes.NewBufferString(`{"main":{"temp":21.5}}`)),
				}
				client.EXPECT().Get("api/v1/weather?lat=53.900000&lon=27.566700&appid=8rfne8&units=metric").Return(resp, nil)
			},
			lat:     53.9,
			lon:     27.5667,
			want:    weather.CurrentWeather{Temperature: 21.5},
			wantErr: false,
		},
		{
			name: "client error",
			prepare: func(client *mock.MockWebClient) {
				client.EXPECT().Get("api/v1/weather?lat=10.000000&lon=20.000000&appid=8rfne8&units=metric").Return(nil, errors.New("some error"))
			},
			lat:     10,
			lon:     20,
			wantErr: true,
		},
		{
			name: "!= 200 status",
			prepare: func(client *mock.MockWebClient) {
				resp := &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(`{"main":{"temp":5}}`)),
				}
				client.EXPECT().Get("api/v1/weather?lat=-1.250000&lon=100.000000&appid=8rfne8&units=metric").Return(resp, nil)
			},
			lat:     -1.25,
			lon:     100,
			wantErr: true,
		},
		{
			name: "invalid json",
			prepare: func(client *mock.MockWebClient) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"main":{"temp":"OGO"}}`)),
				}
				client.EXPECT().Get("api/v1/weather?lat=0.000000&lon=0.000000&appid=8rfne8&units=metric").Return(resp, nil)
			},
			lat:     0,
			lon:     0,
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

func TestOpenWeatherClient_LocationForecast(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpClient := mock.NewMockWebClient(ctrl)
	client, err := clients.NewOpenWeatherClient("8rfne8", "api/v1", httpClient)
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
						`{"list":[{"dt":1713024000,"main":{"temp":11.5}},{"dt":1713034800,"main":{"temp":12.1}}]}`)),
				}
				client.EXPECT().Get("api/v1/forecast?lat=53.900000&lon=27.566700&appid=8rfne8&units=metric").Return(resp, nil)
			},
			lat: 53.9,
			lon: 27.5667,
			want: weather.Forecast{Points: []weather.ForecastPoint{
				{Time: 1713024000, Temperature: 11.5},
				{Time: 1713034800, Temperature: 12.1},
			}},
			wantErr: false,
		},
		{
			name: "client error",
			prepare: func(client *mock.MockWebClient) {
				client.EXPECT().Get("api/v1/forecast?lat=10.000000&lon=20.000000&appid=8rfne8&units=metric").Return(nil, errors.New("network"))
			},
			lat:     10,
			lon:     20,
			wantErr: true,
		},
		{
			name: "!= 200 status",
			prepare: func(client *mock.MockWebClient) {
				resp := &http.Response{
					StatusCode: http.StatusBadGateway,
					Body:       io.NopCloser(bytes.NewBufferString(`{"list":[{"dt":1713024000,"main":{"temp":5}}]}`)),
				}
				client.EXPECT().Get("api/v1/forecast?lat=-1.250000&lon=100.000000&appid=8rfne8&units=metric").Return(resp, nil)
			},
			lat:     -1.25,
			lon:     100,
			wantErr: true,
		},
		{
			name: "invalid json",
			prepare: func(client *mock.MockWebClient) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"list":[{"dt":1713024000,"main":{"temp":"OGO"}}]}`)),
				}
				client.EXPECT().Get("api/v1/forecast?lat=0.000000&lon=0.000000&appid=8rfne8&units=metric").Return(resp, nil)
			},
			lat:     0,
			lon:     0,
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

func TestOpenWeatherClient_CityCurrentTemperature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpClient := mock.NewMockWebClient(ctrl)
	client, err := clients.NewOpenWeatherClient("8rfne8", "api/v1", httpClient)
	require.NoError(t, err)

	tests := []struct {
		name    string
		prepare func(client *mock.MockWebClient)
		city    string
		want    weather.CurrentWeather
		wantErr bool
	}{
		{
			name: "ok",
			prepare: func(client *mock.MockWebClient) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"main":{"temp":21.5}}`)),
				}
				client.EXPECT().Get("api/v1/weather?q=минск&appid=8rfne8&units=metric").Return(resp, nil)
			},
			city:    "минск",
			want:    weather.CurrentWeather{Temperature: 21.5},
			wantErr: false,
		},
		{
			name: "client error",
			prepare: func(client *mock.MockWebClient) {
				client.EXPECT().Get("api/v1/weather?q=лондон&appid=8rfne8&units=metric").Return(nil, errors.New("network"))
			},
			city:    "лондон",
			wantErr: true,
		},
		{
			name: "!= 200 status",
			prepare: func(client *mock.MockWebClient) {
				resp := &http.Response{
					StatusCode: http.StatusBadGateway,
					Body:       io.NopCloser(bytes.NewBufferString(`{"main":{"temp":5}}`)),
				}
				client.EXPECT().Get("api/v1/weather?q=токио&appid=8rfne8&units=metric").Return(resp, nil)
			},
			city:    "токио",
			wantErr: true,
		},
		{
			name: "invalid json",
			prepare: func(client *mock.MockWebClient) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"main":{"temp":"OGO"}}`)),
				}
				client.EXPECT().Get("api/v1/weather?q=шанхай&appid=8rfne8&units=metric").Return(resp, nil)
			},
			city:    "шанхай",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(httpClient)

			got, err := client.CityCurrentTemperature(tt.city)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOpenWeatherClient_CityForecast(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpClient := mock.NewMockWebClient(ctrl)
	client, err := clients.NewOpenWeatherClient("8rfne8", "api/v1", httpClient)
	require.NoError(t, err)

	tests := []struct {
		name    string
		prepare func(client *mock.MockWebClient)
		city    string
		want    weather.Forecast
		wantErr bool
	}{
		{
			name: "ok",
			prepare: func(client *mock.MockWebClient) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(
						`{"list":[{"dt":1713024000,"main":{"temp":8.0}},{"dt":1713034800,"main":{"temp":9.2}}]}`)),
				}
				client.EXPECT().Get("api/v1/forecast?q=минск&appid=8rfne8&units=metric").Return(resp, nil)
			},
			city: "минск",
			want: weather.Forecast{Points: []weather.ForecastPoint{
				{Time: 1713024000, Temperature: 8.0},
				{Time: 1713034800, Temperature: 9.2},
			}},
			wantErr: false,
		},
		{
			name: "client error",
			prepare: func(client *mock.MockWebClient) {
				client.EXPECT().Get("api/v1/forecast?q=лондон&appid=8rfne8&units=metric").Return(nil, errors.New("network"))
			},
			city:    "лондон",
			wantErr: true,
		},
		{
			name: "!= 200 status",
			prepare: func(client *mock.MockWebClient) {
				resp := &http.Response{
					StatusCode: http.StatusBadGateway,
					Body:       io.NopCloser(bytes.NewBufferString(`{"list":[{"dt":1713024000,"main":{"temp":5}}]}`)),
				}
				client.EXPECT().Get("api/v1/forecast?q=токио&appid=8rfne8&units=metric").Return(resp, nil)
			},
			city:    "токио",
			wantErr: true,
		},
		{
			name: "invalid json",
			prepare: func(client *mock.MockWebClient) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"list":[{"dt":1713024000,"main":{"temp":"OGO"}}]}`)),
				}
				client.EXPECT().Get("api/v1/forecast?q=шанхай&appid=8rfne8&units=metric").Return(resp, nil)
			},
			city:    "шанхай",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(httpClient)
			got, err := client.CityForecast(tt.city)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
