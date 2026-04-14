package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/models/weather"
)

type googleCurrentResponse struct {
	Temperature struct {
		Degrees float64 `json:"degrees"`
	} `json:"temperature"`
}

type googleForecastResponse struct {
	ForecastDays []struct {
		Interval struct {
			StartTime string `json:"startTime"`
		} `json:"interval"`
		MaxTemperature struct {
			Degrees float64 `json:"degrees"`
		} `json:"maxTemperature"`
	} `json:"forecastDays"`
}

type GoogleWeatherClient struct {
	apiKey  string
	baseURL string
	client  WebClient
}

func NewGoogleWeatherClient(apiKey, baseURL string, client WebClient) (*GoogleWeatherClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("%w: apiKey is required", ErrInvalidParam)
	}

	if baseURL == "" {
		return nil, fmt.Errorf("%w: baseURL is required", ErrInvalidParam)
	}

	if client == nil {
		return nil, fmt.Errorf("%w: client is required", ErrInvalidParam)
	}

	return &GoogleWeatherClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  client,
	}, nil
}

func (c *GoogleWeatherClient) LocationCurrentTemperature(lat, lon float64) (weather.CurrentWeather, error) {
	url := fmt.Sprintf("%s/currentConditions:lookup?key=%s&location.latitude=%.6f&location.longitude=%.6f",
		c.baseURL, c.apiKey, lat, lon)

	resp, err := c.client.Get(url)
	if err != nil {
		return weather.CurrentWeather{}, fmt.Errorf("failed to call google weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weather.CurrentWeather{}, fmt.Errorf("google weather returned bad status: %d", resp.StatusCode)
	}

	var data googleCurrentResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return weather.CurrentWeather{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return weather.CurrentWeather{Temperature: data.Temperature.Degrees}, nil
}

func (c *GoogleWeatherClient) LocationForecast(lat, lon float64) (weather.Forecast, error) {
	url := fmt.Sprintf("%s/forecast/days:lookup?key=%s&location.latitude=%.6f&location.longitude=%.6f&days=3",
		c.baseURL, c.apiKey, lat, lon)

	resp, err := c.client.Get(url)
	if err != nil {
		return weather.Forecast{}, fmt.Errorf("failed to call google weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weather.Forecast{}, fmt.Errorf("google weather returned bad status: %d", resp.StatusCode)
	}

	var data googleForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return weather.Forecast{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(data.ForecastDays) < 3 {
		return weather.Forecast{}, fmt.Errorf("google weather returned insufficient forecast days: %d", len(data.ForecastDays))
	}

	points := make([]weather.ForecastPoint, 0, 3)

	for i := 0; i < 3; i++ {
		t, err := time.Parse(time.RFC3339, data.ForecastDays[i].Interval.StartTime)
		if err != nil {
			return weather.Forecast{}, fmt.Errorf("invalid startTime: %w", err)
		}

		points = append(points, weather.ForecastPoint{
			Time:        t.Unix(),
			Temperature: data.ForecastDays[i].MaxTemperature.Degrees,
		})
	}

	return weather.Forecast{Points: points}, nil
}

func (c *GoogleWeatherClient) CityCurrentTemperature(city string) (weather.CurrentWeather, error) {
	lat, lon, ok := coordsForCity(city)
	if !ok {
		return weather.CurrentWeather{}, fmt.Errorf("%w: unsupported city %s", ErrInvalidParam, city)
	}

	return c.LocationCurrentTemperature(lat, lon)
}

func (c *GoogleWeatherClient) CityForecast(city string) (weather.Forecast, error) {
	lat, lon, ok := coordsForCity(city)
	if !ok {
		return weather.Forecast{}, fmt.Errorf("%w: unsupported city %s", ErrInvalidParam, city)
	}

	return c.LocationForecast(lat, lon)
}
