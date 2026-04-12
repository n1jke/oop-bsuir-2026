package main

import (
	"log"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-6/api"

	_ "github.com/n1jke/oop-bsuir-2025/laboratory_work-6/api/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Wheather Example API
// @version         1.0
// @BasePath  /api/v1

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	currentWeatherHandler := api.NewCurrentWeatherHandler()

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/weather", currentWeatherHandler.HandleGetCurrentWeather)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run()
}
