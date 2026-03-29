package domain

import "errors"

var (
	ErrEmptyCargoName        = errors.New("cargo name is empty")
	ErrInvalidCargoWeight    = errors.New("cargo weight must be greater than zero")
	ErrInvalidCargoCostPerKg = errors.New("cargo cost per kg must be greater than zero")
)

type CargoInfo struct {
	name      ProductName
	weight    float64
	costPerKg float64
}

type ProductName string

func NewCargoInfo(name ProductName, weight, costPerKg float64) (*CargoInfo, error) {
	if name == "" {
		return nil, ErrEmptyCargoName
	}

	if weight <= 0 {
		return nil, ErrInvalidCargoWeight
	}

	if costPerKg <= 0 {
		return nil, ErrInvalidCargoCostPerKg
	}

	return &CargoInfo{name: name, weight: weight, costPerKg: costPerKg}, nil
}

func (c *CargoInfo) Name() ProductName {
	return c.name
}

func (c *CargoInfo) Weight() float64 {
	return c.weight
}

func (c *CargoInfo) CostPerKg() float64 {
	return c.costPerKg
}

func (c *CargoInfo) CostPerUnit() float64 {
	return c.costPerKg * c.weight
}
