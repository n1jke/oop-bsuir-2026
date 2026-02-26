package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidDistance   = errors.New("distance must not be negative")
	ErrNilTransport      = errors.New("transport must not be nil")
	ErrEmptyOrderContent = errors.New("order content must not be empty")
)

type Order struct {
	id        uuid.UUID
	dist      float64
	transport Transport
	content   []ProductBatch
}

type ProductBatch struct {
	CargoInfo
	count uint
}

func NewOrder(dist float64, transport Transport, content []ProductBatch) (*Order, error) {
	if dist < 0 {
		return nil, ErrInvalidDistance
	}

	if transport == nil {
		return nil, ErrNilTransport
	}

	if len(content) == 0 {
		return nil, ErrEmptyOrderContent
	}

	return &Order{
		id:        uuid.New(),
		dist:      dist,
		transport: transport,
		content:   content,
	}, nil
}

func (o *Order) ID() uuid.UUID {
	return o.id
}

func (o *Order) CalculateCost() float64 {
	total := 0.0

	for i := range o.content {
		total += o.content[i].CostForBatch()
	}

	total += o.transport.DeliveryRate() * o.dist

	return total
}

func (o *Order) EstimateDuration() time.Duration {
	speed := o.transport.Speed()
	if speed <= 0 {
		return 0
	}

	hours := o.dist / speed

	return time.Duration(hours * float64(time.Hour))
}

func NewProductBatch(item CargoInfo, count uint) *ProductBatch {
	return &ProductBatch{CargoInfo: item, count: count}
}

func (p *ProductBatch) CostForBatch() float64 {
	return float64(p.count) * p.CostPerUnit()
}

func (p *ProductBatch) TotalWeight() float64 {
	return float64(p.count) * p.Weight()
}
