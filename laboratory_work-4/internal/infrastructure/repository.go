package infrastructure

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-4/internal/domain"
)

const (
	recordType int = iota
	name
	transportMode
	weightPerKg
	deliveryCost
	fuelRate
	speed
)

type CsvRepository struct {
	logger *slog.Logger
	path   string
}

func NewCsvRepository(logger *slog.Logger, path string) *CsvRepository {
	return &CsvRepository{
		logger: logger,
		path:   path,
	}
}

func (c CsvRepository) LoadCargoInfo() ([]domain.CargoInfo, error) {
	records, err := c.readRecords(c.path)
	if err != nil {
		return nil, err
	}

	output := make([]domain.CargoInfo, 0, len(records))

	for i := 1; i < len(records); i++ {
		if len(records[i]) <= speed {
			c.logger.Error("invalid csv row length", "row", i, "len", len(records[i]))
			continue
		}
		if records[i][recordType] != "cargo" {
			continue
		}

		cName := strings.TrimSpace(records[i][name])

		weight, err := strconv.ParseFloat(records[i][weightPerKg], 64)
		if err != nil {
			c.logger.Error("error parsing cargo weight", "row", i, "error", err)
			continue
		}

		cost, err := strconv.ParseFloat(records[i][deliveryCost], 64)
		if err != nil {
			c.logger.Error("error parsing cargo cost", "row", i, "error", err)
			continue
		}

		info, err := domain.NewCargoInfo(domain.ProductName(cName), weight, cost)
		if err != nil {
			c.logger.Error("error creating cargo info", "row", i, "error", err)
			continue
		}

		output = append(output, *info)
	}

	return output, nil
}

func (c CsvRepository) LoadTransportInfo() ([]domain.TransportInfo, error) {
	records, err := c.readRecords(c.path)
	if err != nil {
		return nil, err
	}

	output := make([]domain.TransportInfo, 0, len(records))

	for i := range records {
		if len(records[i]) <= speed {
			c.logger.Error("invalid csv row length", "row", i, "len", len(records[i]))
			continue
		}
		if records[i][recordType] != "transport" {
			continue
		}

		tName := strings.TrimSpace(records[i][name])

		mode, err := parseTransportMode(records[i][transportMode])

		rate, err := strconv.ParseFloat(records[i][fuelRate], 64)
		if err != nil {
			c.logger.Error("error parsing transport rate", "row", i, "error", err)
			continue
		}

		transportSpeed, err := strconv.ParseFloat(records[i][speed], 64)
		if err != nil {
			c.logger.Error("error parsing transport speed", "row", i, "error", err)
			continue
		}

		info, err := domain.NewTransportInfo(domain.TransportType(tName), domain.TransportMode(mode), rate, transportSpeed)
		if err != nil {
			c.logger.Error("error creating transport info", "row", i, "error", err)
			continue
		}

		output = append(output, *info)
	}

	return output, nil
}

func (c CsvRepository) readRecords(path string) ([][]string, error) {
	csvfile, err := os.Open(path)
	if err != nil {
		c.logger.Error("Error while opening configuration file", "error", err)
		return nil, err
	}
	defer func() {
		if err := csvfile.Close(); err != nil {
			c.logger.Error("Error while closing configuration file", "error", err)
		}
	}()

	reader := csv.NewReader(csvfile)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		c.logger.Error("Error while reading configuration file", "error", err)
		return nil, err
	}

	return records, nil
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
