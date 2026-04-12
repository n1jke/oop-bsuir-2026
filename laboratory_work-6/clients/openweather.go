package clients

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
)

type openWeatherResponse struct {
	Main struct {
		Temp decimal.Decimal `json:"temp"`
	} `json:"main"`
}

type OpenWeatherClient struct {
	apiKey  string
	baseURL string
}

func NewOpenWeatherClient(apiKey string, baseURL string) *OpenWeatherClient {
	return &OpenWeatherClient{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// Implementation of WeatherDataClient
func (c *OpenWeatherClient) LocationCurrentTemperature(lat decimal.Decimal, lon decimal.Decimal) (decimal.Decimal, error) {
	url := fmt.Sprintf("%s?lat=%s&lon=%s&appid=%s&units=metric",
		c.baseURL, lat.String(), lon.String(), c.apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to call openweather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return decimal.Zero, fmt.Errorf("openweather returned bad status: %d", resp.StatusCode)
	}

	var data openWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return decimal.Zero, fmt.Errorf("failed to decode response: %w", err)
	}

	return data.Main.Temp, nil
}
