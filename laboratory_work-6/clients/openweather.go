package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type openWeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

type OpenWeatherClient struct {
	apiKey  string
	baseURL string
	client  WebClient
}

func NewOpenWeatherClient(apiKey, baseURL string, client WebClient) *OpenWeatherClient {
	return &OpenWeatherClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  client,
	}
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
