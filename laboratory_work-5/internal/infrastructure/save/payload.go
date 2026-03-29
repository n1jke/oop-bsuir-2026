package save

import (
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/application"
)

type ExportPayload struct {
	Order   *OrderPayload  `json:"order,omitempty" yaml:"order,omitempty"`
	Cost    float64        `json:"cost" yaml:"cost"`
	Options []QuotePayload `json:"options,omitempty" yaml:"options,omitempty"`
}

type OrderPayload struct {
	ID            string  `json:"id" yaml:"id"`
	Distance      float64 `json:"distance" yaml:"distance"`
	TransportName string  `json:"transport_name,omitempty" yaml:"transport_name,omitempty"`
	TransportMode string  `json:"transport_mode,omitempty" yaml:"transport_mode,omitempty"`
	Duration      string  `json:"duration" yaml:"duration"`
}

type QuotePayload struct {
	TransportName string  `json:"transport_name" yaml:"transport_name"`
	TransportMode string  `json:"transport_mode" yaml:"transport_mode"`
	Cost          float64 `json:"cost" yaml:"cost"`
	Duration      string  `json:"duration" yaml:"duration"`
}

func mapServiceResponse(resp *application.ServiceResponse) ExportPayload {
	out := ExportPayload{
		Cost:    resp.Cost,
		Options: make([]QuotePayload, 0, len(resp.Options)),
	}

	for i := range resp.Options {
		q := resp.Options[i]
		out.Options = append(out.Options, QuotePayload{
			TransportName: string(q.Transport.Name()),
			TransportMode: string(q.Transport.Mode()),
			Cost:          q.Cost,
			Duration:      q.Duration.String(),
		})
	}

	if resp.Order != nil {
		orderPayload := OrderPayload{
			ID:       resp.Order.ID().String(),
			Distance: resp.Order.Distance(),
			Duration: resp.Order.EstimateDuration().String(),
		}

		if tr := resp.Order.Transport(); tr != nil {
			orderPayload.TransportName = string(tr.Name())
			orderPayload.TransportMode = string(tr.Mode())
		}

		out.Order = &orderPayload
	}

	return out
}
