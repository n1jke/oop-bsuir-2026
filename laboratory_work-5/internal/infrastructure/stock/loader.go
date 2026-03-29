package stock

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/application"
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/domain"
)

var ErrInvalidFileType = errors.New("provided invalid file type for stock loader")

func NewStockLoader(logger *slog.Logger, path, loaderType string) (application.StockLoader, error) {
	switch loaderType {
	case "csv":
		return NewCsvLoader(logger, path), nil
	case "json":
		return NewJSONLoader(logger, path), nil
	case "xml":
		return NewXMLLoader(logger, path), nil
	default:
		return nil, ErrInvalidFileType
	}
}

func parseTransportMode(raw string) (domain.TransportMode, error) {
	switch raw {
	case "Земля", "land":
		return domain.LandTransport, nil
	case "Вода", "water":
		return domain.WaterTransport, nil
	case "Воздух", "air":
		return domain.AirTransport, nil
	default:
		return "", fmt.Errorf("unknown transport mode: %q", raw)
	}
}
