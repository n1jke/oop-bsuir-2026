package api

import (
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-6/clients"
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-6/controllers"
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-6/models/weather"
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-6/shared/responses"
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-6/shared/utils"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type WeatherHandler struct {
	Controller controllers.CurrentWeatherController[*clients.OpenWeatherClient]
}

func NewCurrentWeatherHandler() *WeatherHandler {
	return &WeatherHandler{Controller: *controllers.NewCurrentWeatherController(
		clients.NewOpenWeatherClient(utils.GetEnv("OPENWEATHER_API_KEY", ""), utils.GetEnv("OPENWEATHER_BASE_URL", "")))}
}

// GetCurrentWeather godoc
// @Summary      Get Current Weather
// @Description  Returns current weather for given coordinates
// @Tags         weather
// @Produce      json
// @Param        lat   query     string  true  "Latitude"    default(18.300231990440125)
// @Param        lon   query     string  true  "Longitude"   default(-64.8251590359234)
// @Success      200   {object}  responses.SuccessResponse[weather.CurrentWeather]
// @Failure      400   {object}  responses.StatusResponse
// @Failure      500   {object}  responses.StatusResponse
// @Router       /weather [get]
func (h *WeatherHandler) HandleGetCurrentWeather(c *gin.Context) {
	lat, errLat := decimal.NewFromString(c.Query("lat"))
	lon, errLon := decimal.NewFromString(c.Query("lon"))

	if errLat != nil || errLon != nil {
		c.JSON(400, responses.StatusResponse{Code: 400, Message: "invalid coordinates"})
		return
	}

	result, err := h.Controller.GetCurrentWeather(lat, lon)
	if err != nil {

		c.JSON(500, responses.StatusResponse{Code: 500, Message: err.Error()})
		return
	}
	c.JSON(200, responses.SuccessResponse[weather.CurrentWeather]{Code: 200, Message: "Success", Data: result})
}
