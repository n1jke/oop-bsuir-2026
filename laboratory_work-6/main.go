package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/api"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/shared/responses"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-6/shared/utils"
)

func main() {
	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	currentWeatherHandler, err := api.NewCurrentWeatherHandler(cfg.OpenWeatherKey, cfg.OpenWeatherURL)
	if err != nil {
		log.Fatalf("Error creating weather handler: %v", err)
	}

	r := gin.Default()
	api.RegisterHandlersWithOptions(r, currentWeatherHandler, api.GinServerOptions{
		BaseURL: "/api/v1",
		ErrorHandler: func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, responses.StatusResponse{Code: statusCode, Message: err.Error()})
		},
	})

	_ = r.Run("localhost:8080")
}
