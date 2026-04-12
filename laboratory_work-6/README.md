# Weather Forecast Example Project

## Docs
[swagger](http://localhost:8080/swagger/index.html)

## Base commands
- docs gen: `swag init -g main.go -o ./api/docs --parseDependency`
> just regenerate for update

- run: `go run main.go`

## Environment variables
`.env` example:
```
OPENWEATHER_API_KEY=xxxxxxxxxxxxxxxxx
OPENWEATHER_BASE_URL=https://api.openweathermap.org/data/2.5/weather
```
