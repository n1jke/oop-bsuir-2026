package clients

import "github.com/shopspring/decimal"

type WeatherDataClient interface {
	LocationCurrentTemperature(lat decimal.Decimal, lon decimal.Decimal) (temperature decimal.Decimal, err error)
}
