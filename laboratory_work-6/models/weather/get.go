package weather

import "github.com/shopspring/decimal"

type CurrentWeather struct {
	Temperature decimal.Decimal `json:"temperature"`
}
