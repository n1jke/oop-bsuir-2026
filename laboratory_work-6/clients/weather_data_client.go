package clients

import "net/http"

type WeatherDataClient interface {
	LocationCurrentTemperature(lat, lon float64) (temperature float64, err error)
}

type WebClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (resp *http.Response, err error)
}
