package stock

import (
	"encoding/csv"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/domain"
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

type CsvLoader struct {
	logger *slog.Logger
	path   string
}

func NewCsvLoader(logger *slog.Logger, path string) *CsvLoader {
	return &CsvLoader{
		logger: logger,
		path:   path,
	}
}

func (cl CsvLoader) LoadCargoInfo() ([]domain.CargoInfo, error) {
	records, err := cl.readRecords(cl.path)
	if err != nil {
		return nil, err
	}

	output := make([]domain.CargoInfo, 0, len(records))

	for i := 1; i < len(records); i++ {
		if len(records[i]) <= speed {
			cl.logger.Error("invalid csv row length", "row", i, "len", len(records[i]))
			continue
		}

		if records[i][recordType] != "cargo" {
			continue
		}

		cName := strings.TrimSpace(records[i][name])

		weight, err := strconv.ParseFloat(records[i][weightPerKg], 64)
		if err != nil {
			cl.logger.Error("error parsing cargo weight", "row", i, "error", err)
			continue
		}

		cost, err := strconv.ParseFloat(records[i][deliveryCost], 64)
		if err != nil {
			cl.logger.Error("error parsing cargo cost", "row", i, "error", err)
			continue
		}

		info, err := domain.NewCargoInfo(domain.ProductName(cName), weight, cost)
		if err != nil {
			cl.logger.Error("error creating cargo info", "row", i, "error", err)
			continue
		}

		output = append(output, *info)
	}

	return output, nil
}

func (cl CsvLoader) LoadTransportInfo() ([]domain.TransportInfo, error) {
	records, err := cl.readRecords(cl.path)
	if err != nil {
		return nil, err
	}

	output := make([]domain.TransportInfo, 0, len(records))

	for i := range records {
		if len(records[i]) <= speed {
			cl.logger.Error("invalid csv row length", "row", i, "len", len(records[i]))
			continue
		}

		if records[i][recordType] != "transport" {
			continue
		}

		tName := strings.TrimSpace(records[i][name])

		mode, err := parseTransportMode(records[i][transportMode])
		if err != nil {
			cl.logger.Error("error parsing transport mode", "row", i, "error", err)
			continue
		}

		rate, err := strconv.ParseFloat(records[i][fuelRate], 64)
		if err != nil {
			cl.logger.Error("error parsing transport rate", "row", i, "error", err)
			continue
		}

		transportSpeed, err := strconv.ParseFloat(records[i][speed], 64)
		if err != nil {
			cl.logger.Error("error parsing transport speed", "row", i, "error", err)
			continue
		}

		info, err := domain.NewTransportInfo(domain.TransportType(tName), mode, rate, transportSpeed)
		if err != nil {
			cl.logger.Error("error creating transport info", "row", i, "error", err)
			continue
		}

		output = append(output, *info)
	}

	return output, nil
}

func (cl CsvLoader) readRecords(path string) ([][]string, error) {
	csvfile, err := os.Open(path)
	if err != nil {
		cl.logger.Error("Error while opening configuration file", "error", err)
		return nil, err
	}

	defer func() {
		if err := csvfile.Close(); err != nil {
			cl.logger.Error("Error while closing configuration file", "error", err)
		}
	}()

	reader := csv.NewReader(csvfile)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		cl.logger.Error("Error while reading configuration file", "error", err)
		return nil, err
	}

	return records, nil
}
