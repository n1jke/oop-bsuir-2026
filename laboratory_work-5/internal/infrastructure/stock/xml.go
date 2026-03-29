package stock

import (
	"encoding/xml"
	"log/slog"
	"os"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/domain"
)

type XMLLoader struct {
	logger *slog.Logger
	path   string
}

func (xl XMLLoader) LoadCargoInfo() ([]domain.CargoInfo, error) {
	cfg, err := xl.readConfig()
	if err != nil {
		return nil, err
	}

	output := make([]domain.CargoInfo, 0, len(cfg.Cargo))
	for i := range cfg.Cargo {
		item := cfg.Cargo[i]

		info, err := domain.NewCargoInfo(domain.ProductName(item.Name), item.Weight, item.CostPerKg)
		if err != nil {
			xl.logger.Error("error creating cargo info", "row", i, "error", err)
			continue
		}

		output = append(output, *info)
	}

	return output, nil
}

func (xl XMLLoader) LoadTransportInfo() ([]domain.TransportInfo, error) {
	cfg, err := xl.readConfig()
	if err != nil {
		return nil, err
	}

	output := make([]domain.TransportInfo, 0, len(cfg.Transport))
	for i := range cfg.Transport {
		item := cfg.Transport[i]

		mode, err := parseTransportMode(item.Mode)
		if err != nil {
			xl.logger.Error("error parsing transport mode", "row", i, "error", err)
			continue
		}

		info, err := domain.NewTransportInfo(domain.TransportType(item.Name), mode, item.RatePerKm, item.Speed)
		if err != nil {
			xl.logger.Error("error creating transport info", "row", i, "error", err)
			continue
		}

		output = append(output, *info)
	}

	return output, nil
}

func NewXMLLoader(logger *slog.Logger, path string) *XMLLoader {
	return &XMLLoader{
		logger: logger,
		path:   path,
	}
}

type xmlConfig struct {
	XMLName   xml.Name       `xml:"config"`
	Cargo     []xmlCargo     `xml:"cargo"`
	Transport []xmlTransport `xml:"transport"`
}

type xmlCargo struct {
	Name      string  `xml:"name"`
	Weight    float64 `xml:"weight"`
	CostPerKg float64 `xml:"cost_per_kg"`
}

type xmlTransport struct {
	Name      string  `xml:"name"`
	Mode      string  `xml:"mode"`
	RatePerKm float64 `xml:"rate_per_km"`
	Speed     float64 `xml:"speed"`
}

func (xl XMLLoader) readConfig() (*xmlConfig, error) {
	data, err := os.ReadFile(xl.path)
	if err != nil {
		xl.logger.Error("error reading xml config file", "error", err)
		return nil, err
	}

	var cfg xmlConfig
	if err := xml.Unmarshal(data, &cfg); err != nil {
		xl.logger.Error("error parsing xml config file", "error", err)
		return nil, err
	}

	return &cfg, nil
}
