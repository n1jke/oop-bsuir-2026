package clients

import "strings"

type cityCoord struct {
	lat float64
	lon float64
}

var cityCoords = map[string]cityCoord{
	"минск":   {lat: 53.9000, lon: 27.5667},
	"лондон":  {lat: 51.5074, lon: -0.1278},
	"шанхай":  {lat: 31.2304, lon: 121.4737},
	"варшава": {lat: 52.2297, lon: 21.0122},
}

func coordsForCity(city string) (float64, float64, bool) {
	key := strings.ToLower(strings.TrimSpace(city))

	coord, ok := cityCoords[key]
	if !ok {
		return 0, 0, false
	}

	return coord.lat, coord.lon, true
}
