package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidDistance   = errors.New("distance must not be negative")
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

	// todo: maybe temp solution(case: client dont choose a transport & we create list of possible)
	// if transport == nil {
	// 	return nil, ErrNilTransport
	// }

	if len(content) == 0 {
		return nil, ErrEmptyOrderContent
	}

	contentCopy := make([]ProductBatch, len(content))
	copy(contentCopy, content)

	return &Order{
		id:        uuid.New(),
		dist:      dist,
		transport: transport,
		content:   contentCopy,
	}, nil
}

func (o *Order) ID() uuid.UUID {
	return o.id
}

func (o *Order) Transport() Transport {
	return o.transport
}

func (o *Order) Distance() float64 {
	return o.dist
}

func (o *Order) CalculateCost() float64 {
	total := 0.0

	for i := range o.content {
		total += o.content[i].CostForBatch()
	}

	if o.transport != nil {
		total += o.transport.CalculateCost(o.dist)
	}

	return total
}

func (o *Order) EstimateDuration() time.Duration {
	if o.transport == nil {
		return time.Second * 0
	}

	return time.Duration(o.transport.CalculateDeliveryTime(o.dist) * float64(time.Hour))
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
