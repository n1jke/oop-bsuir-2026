package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/controllers"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/models/weather"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/shared/responses"
)

type WeatherHandler struct {
	Controller controllers.CurrentWeatherController[controllers.WeatherDataClient]
}

func NewCurrentWeatherHandler(client controllers.WeatherDataClient) (*WeatherHandler, error) {
	return &WeatherHandler{
		Controller: *controllers.NewCurrentWeatherController(client),
	}, nil
}

// todo: can update with concurretly logic of crawling weather for cities and coordinates

func (h *WeatherHandler) GetWeatherCity(c *gin.Context, params GetWeatherCityParams) {
	resp := make([]weather.CurrentWeather, 0, len(params.City))

	for i := range params.City {
		result, err := h.Controller.GetWeatherCity(params.City[i])
		if err != nil {
			statusCode := 500
			if errors.Is(err, controllers.ErrUnsupportedCity) {
				statusCode = 400
			}

			c.JSON(statusCode, responses.StatusResponse{Code: statusCode, Message: err.Error()})

			return
		}

		resp = append(resp, result)
	}

	c.JSON(200, responses.SuccessResponse[[]weather.CurrentWeather]{Code: 200, Message: "Success", Data: resp})
}

func (h *WeatherHandler) GetWeatherCoordinate(c *gin.Context, params GetWeatherCoordinateParams) {
	resp := make([]weather.CurrentWeather, 0, len(params.Coord))

	for i := range params.Coord {
		lat, lon, err := parseCoord(params.Coord[i])
		if err != nil {
			c.JSON(400, responses.StatusResponse{Code: 400, Message: err.Error()})
			return
		}

		result, err := h.Controller.GetWeatherCoordinates(lat, lon)
		if err != nil {
			c.JSON(500, responses.StatusResponse{Code: 500, Message: err.Error()})
			return
		}

		resp = append(resp, result)
	}

	c.JSON(200, responses.SuccessResponse[[]weather.CurrentWeather]{Code: 200, Message: "Success", Data: resp})
}

func (h *WeatherHandler) GetForecastCity(c *gin.Context, params GetForecastCityParams) {
	resp := make([]weather.Forecast, 0, len(params.City))

	for i := range params.City {
		result, err := h.Controller.GetForecastCity(params.City[i])
		if err != nil {
			statusCode := 500
			if errors.Is(err, controllers.ErrUnsupportedCity) {
				statusCode = 400
			}

			c.JSON(statusCode, responses.StatusResponse{Code: statusCode, Message: err.Error()})

			return
		}

		resp = append(resp, result)
	}

	c.JSON(200, responses.SuccessResponse[[]weather.Forecast]{Code: 200, Message: "Success", Data: resp})
}

func (h *WeatherHandler) GetForecastCoordinate(c *gin.Context, params GetForecastCoordinateParams) {
	resp := make([]weather.Forecast, 0, len(params.Coord))

	for i := range params.Coord {
		lat, lon, err := parseCoord(params.Coord[i])
		if err != nil {
			c.JSON(400, responses.StatusResponse{Code: 400, Message: err.Error()})
			return
		}

		result, err := h.Controller.GetForecastCoordinates(lat, lon)
		if err != nil {
			c.JSON(500, responses.StatusResponse{Code: 500, Message: err.Error()})
			return
		}

		resp = append(resp, result)
	}

	c.JSON(200, responses.SuccessResponse[[]weather.Forecast]{Code: 200, Message: "Success", Data: resp})
}

func parseCoord(value string) (float64, float64, error) {
	parts := strings.Split(value, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid coordinate format: %s", value)
	}

	lat, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid latitude: %s", parts[0])
	}

	lon, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid longitude: %s", parts[1])
	}

	return lat, lon, nil
}
