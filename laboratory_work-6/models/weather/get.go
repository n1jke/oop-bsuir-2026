package weather

type CurrentWeather struct {
	Temperature float64
}

type Forecast struct {
	Points []ForecastPoint
}

type ForecastPoint struct {
	Time        int64
	Temperature float64
}
