package application

import (
	"time"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/domain"
)

type ClientResponse struct {
	Transport domain.Transport
	Dist      float64
	Content   []domain.ProductBatch
}

type ServiceResponse struct {
	Order   *domain.Order
	Cost    float64
	Options []Quote
}

func (c ClientResponse) ToOrder() (*domain.Order, error) {
	return domain.NewOrder(c.Dist, c.Transport, c.Content)
}

type Quote struct {
	Transport domain.TransportInfo
	Cost      float64
	Duration  time.Duration
}
