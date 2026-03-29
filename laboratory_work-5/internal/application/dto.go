package application

import (
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/domain"
)

var ErrInvalidSortField = errors.New("invalid sort field")

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

type QuoteCmp func(a, b Quote) int

func CombineCmps(cmps ...QuoteCmp) QuoteCmp {
	return func(a, b Quote) int {
		for _, cmp := range cmps {
			if res := cmp(a, b); res != 0 {
				return res
			}
		}

		return 0
	}
}

func (sr *ServiceResponse) SortFields(fields []string) error {
	if sr == nil || len(sr.Options) == 0 || len(fields) == 0 {
		return nil
	}

	cmps := make([]QuoteCmp, 0, len(fields))
	for i := range fields {
		switch strings.ToLower(strings.TrimSpace(fields[i])) {
		case "cost":
			cmps = append(cmps, compareByCost)
		case "duration":
			cmps = append(cmps, compareByDuration)
		case "transport", "transportname", "transport_name":
			cmps = append(cmps, compareByTransportName)
		case "":
			continue
		default:
			return ErrInvalidSortField
		}
	}

	slices.SortStableFunc(sr.Options, CombineCmps(cmps...))

	return nil
}
