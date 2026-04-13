package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/clients"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/controllers"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/models/weather"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/shared/responses"
)

type WeatherHandler struct {
	Controller controllers.CurrentWeatherController[*clients.OpenWeatherClient]
}

func NewCurrentWeatherHandler(key, url string) (*WeatherHandler, error) {
	client, err := clients.NewOpenWeatherClient(key, url, http.DefaultClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create weather client: %w", err)
	}

	return &WeatherHandler{
		Controller: *controllers.NewCurrentWeatherController(client),
	}, nil
}

func (h *WeatherHandler) GetWeather(c *gin.Context, params GetWeatherParams) {
	result, err := h.Controller.GetWeather(params.Lat, params.Lon)
	if err != nil {
		c.JSON(500, responses.StatusResponse{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(200, responses.SuccessResponse[weather.CurrentWeather]{Code: 200, Message: "Success", Data: result})
}
