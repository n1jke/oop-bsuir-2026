package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/models/weather"
)

var ErrInvalidParam = errors.New("invalid param")

type openWeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	List []struct {
		Dt   int64 `json:"dt"`
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
	} `json:"list"`
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
func (c *OpenWeatherClient) LocationCurrentTemperature(lat, lon float64) (weather.CurrentWeather, error) {
	url := fmt.Sprintf("%s/weather?lat=%.6f&lon=%.6f&appid=%s&units=metric", c.baseURL, lat, lon, c.apiKey)

	resp, err := c.client.Get(url)
	if err != nil {
		return weather.CurrentWeather{}, fmt.Errorf("failed to call openweather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weather.CurrentWeather{}, fmt.Errorf("openweather returned bad status: %d", resp.StatusCode)
	}

	var data openWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return weather.CurrentWeather{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return weather.CurrentWeather{Temperature: data.Main.Temp}, nil
}

func (c *OpenWeatherClient) LocationForecast(lat, lon float64) (forecast weather.Forecast, err error) {
	url := fmt.Sprintf("%s/forecast?lat=%.6f&lon=%.6f&appid=%s&units=metric", c.baseURL, lat, lon, c.apiKey)

	resp, err := c.client.Get(url)
	if err != nil {
		return weather.Forecast{}, fmt.Errorf("failed to call openweather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weather.Forecast{}, fmt.Errorf("openweather returned bad status: %d", resp.StatusCode)
	}

	var data openWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return weather.Forecast{}, fmt.Errorf("failed to decode response: %w", err)
	}

	// map to domain model
	points := make([]weather.ForecastPoint, 0, len(data.List))
	for i := range data.List {
		points = append(points, weather.ForecastPoint{
			Time:        data.List[i].Dt,
			Temperature: data.List[i].Main.Temp,
		})
	}

	return weather.Forecast{Points: points}, nil
}

func (c *OpenWeatherClient) CityCurrentTemperature(city string) (weather.CurrentWeather, error) {
	url := fmt.Sprintf("%s/weather?q=%s&appid=%s&units=metric", c.baseURL, city, c.apiKey)

	resp, err := c.client.Get(url)
	if err != nil {
		return weather.CurrentWeather{}, fmt.Errorf("failed to call openweather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weather.CurrentWeather{}, fmt.Errorf("openweather returned bad status: %d", resp.StatusCode)
	}

	var data openWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return weather.CurrentWeather{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return weather.CurrentWeather{Temperature: data.Main.Temp}, nil
}

func (c *OpenWeatherClient) CityForecast(city string) (forecast weather.Forecast, err error) {
	url := fmt.Sprintf("%s/forecast?q=%s&appid=%s&units=metric", c.baseURL, city, c.apiKey)

	resp, err := c.client.Get(url)
	if err != nil {
		return weather.Forecast{}, fmt.Errorf("failed to call openweather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weather.Forecast{}, fmt.Errorf("openweather returned bad status: %d", resp.StatusCode)
	}

	var data openWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return weather.Forecast{}, fmt.Errorf("failed to decode response: %w", err)
	}

	// map to domain model
	points := make([]weather.ForecastPoint, 0, len(data.List))
	for i := range data.List {
		points = append(points, weather.ForecastPoint{
			Time:        data.List[i].Dt,
			Temperature: data.List[i].Main.Temp,
		})
	}

	return weather.Forecast{Points: points}, nil
}
