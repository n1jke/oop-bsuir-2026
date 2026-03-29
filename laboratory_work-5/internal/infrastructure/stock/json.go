package stock

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/domain"
)

type JSONLoader struct {
	logger *slog.Logger
	path   string
}

func (jl JSONLoader) LoadCargoInfo() ([]domain.CargoInfo, error) {
	cfg, err := jl.readConfig()
	if err != nil {
		return nil, err
	}

	output := make([]domain.CargoInfo, 0, len(cfg.Cargo))
	for i := range cfg.Cargo {
		item := cfg.Cargo[i]

		info, err := domain.NewCargoInfo(domain.ProductName(item.Name), item.Weight, item.CostPerKg)
		if err != nil {
			jl.logger.Error("error creating cargo info", "row", i, "error", err)
			continue
		}

		output = append(output, *info)
	}

	return output, nil
}

func (jl JSONLoader) LoadTransportInfo() ([]domain.TransportInfo, error) {
	cfg, err := jl.readConfig()
	if err != nil {
		return nil, err
	}

	output := make([]domain.TransportInfo, 0, len(cfg.Transport))
	for i := range cfg.Transport {
		item := cfg.Transport[i]

		mode, err := parseTransportMode(item.Mode)
		if err != nil {
			jl.logger.Error("error parsing transport mode", "row", i, "error", err)
			continue
		}

		info, err := domain.NewTransportInfo(domain.TransportType(item.Name), mode, item.RatePerKm, item.Speed)
		if err != nil {
			jl.logger.Error("error creating transport info", "row", i, "error", err)
			continue
		}

		output = append(output, *info)
	}

	return output, nil
}

func NewJSONLoader(logger *slog.Logger, path string) *JSONLoader {
	return &JSONLoader{
		logger: logger,
		path:   path,
	}
}

type jsonConfig struct {
	Cargo     []jsonCargo     `json:"cargo"`
	Transport []jsonTransport `json:"transport"`
}

type jsonCargo struct {
	Name      string  `json:"name"`
	Weight    float64 `json:"weight"`
	CostPerKg float64 `json:"cost_per_kg"`
}

type jsonTransport struct {
	Name      string  `json:"name"`
	Mode      string  `json:"mode"`
	RatePerKm float64 `json:"rate_per_km"`
	Speed     float64 `json:"speed"`
}

func (jl JSONLoader) readConfig() (*jsonConfig, error) {
	data, err := os.ReadFile(jl.path)
	if err != nil {
		jl.logger.Error("error reading json config file", "error", err)
		return nil, err
	}

	var cfg jsonConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		jl.logger.Error("error parsing json config file", "error", err)
		return nil, err
	}

	return &cfg, nil
}
