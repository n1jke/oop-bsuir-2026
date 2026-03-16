package domain

import "errors"

var (
	ErrEmptyTransportType    = errors.New("transport type is empty")
	ErrEmptyTransportMode    = errors.New("transport mode is empty")
	ErrInvalidRatePerKm      = errors.New("transport rate per km must be greater than zero")
	ErrInvalidTransportSpeed = errors.New("transport speed must be greater than zero")
)

type Transport interface {
	Name() TransportType
	Mode() TransportMode
	DeliveryRate() float64
	Speed() float64
}

type TransportInfo struct {
	transportType TransportType
	mode          TransportMode
	deliveryRate  float64
	speed         float64
}

type (
	TransportType string
	TransportMode string
)

const (
	LandTransport  TransportMode = "land"
	WaterTransport TransportMode = "water"
	AirTransport   TransportMode = "air"
)

func NewTransportInfo(transportType TransportType, mode TransportMode, ratePerKm, speed float64) (*TransportInfo, error) {
	if transportType == "" {
		return nil, ErrEmptyTransportType
	}

	if mode == "" {
		return nil, ErrEmptyTransportMode
	}

	if ratePerKm <= 0 {
		return nil, ErrInvalidRatePerKm
	}

	if speed <= 0 {
		return nil, ErrInvalidTransportSpeed
	}

	return &TransportInfo{
		transportType: transportType,
		mode:          mode,
		deliveryRate:  ratePerKm,
		speed:         speed,
	}, nil
}

func (t *TransportInfo) Name() TransportType {
	return t.transportType
}

func (t *TransportInfo) Mode() TransportMode {
	return t.mode
}

func (t *TransportInfo) DeliveryRate() float64 {
	return t.deliveryRate
}

func (t *TransportInfo) Speed() float64 {
	return t.speed
}
