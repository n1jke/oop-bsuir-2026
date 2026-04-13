package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var ErrInvalidParam = errors.New("invalid param")

type openWeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

//go:generate mockgen -source=openweather.go -destination=mock/mock.go -package=mock
type WebClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (resp *http.Response, err error)
}

type OpenWeatherClient struct {
	apiKey  string
	baseURL string
	client  WebClient
}

func NewOpenWeatherClient(apiKey, baseURL string, client WebClient) (*OpenWeatherClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("%w: apiKey is required", ErrInvalidParam)
	}

	if baseURL == "" {
		return nil, fmt.Errorf("%w: baseURL is required", ErrInvalidParam)
	}

	if client == nil {
		return nil, fmt.Errorf("%w: client is required", ErrInvalidParam)
	}

	return &OpenWeatherClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  client,
	}, nil
}

// Implementation of WeatherDataClient.
func (c *OpenWeatherClient) LocationCurrentTemperature(lat, lon float64) (float64, error) {
	url := fmt.Sprintf("%s?lat=%.6f&lon=%.6f&appid=%s&units=metric", c.baseURL, lat, lon, c.apiKey)

	resp, err := c.client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to call openweather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("openweather returned bad status: %d", resp.StatusCode)
	}

	var data openWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return data.Main.Temp, nil
}
