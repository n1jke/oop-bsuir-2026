package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/domain"
)

func TestOrderCalculateCost(t *testing.T) {
	t.Parallel()

	plane, err := domain.NewTransportInfo("plane", domain.AirTransport, 150, 850)
	require.NoError(t, err)

	tanker, err := domain.NewTransportInfo("tanker", domain.WaterTransport, 2, 35)
	require.NoError(t, err)

	perishables, err := domain.NewCargoInfo("perishables", 10, 100)
	require.NoError(t, err)

	equipment, err := domain.NewCargoInfo("equipment", 120, 15)
	require.NoError(t, err)

	tests := []struct {
		name      string
		dist      float64
		transport domain.Transport
		content   []domain.ProductBatch
		want      float64
	}{
		{
			name:      "only cargo component : distance is zero",
			dist:      0,
			transport: tanker,
			content: []domain.ProductBatch{
				*domain.NewProductBatch(*equipment, 1),
			},
			want: 1 * (120 * 15),
		},
		{
			name:      "cargo with air delivery",
			dist:      250,
			transport: plane,
			content: []domain.ProductBatch{
				*domain.NewProductBatch(*perishables, 2),
				*domain.NewProductBatch(*equipment, 1),
			},
			want: 2*(10*100) + 1*(120*15) + 250*150,
		},
		{
			name:      "cargo with water delivery",
			dist:      200,
			transport: tanker,
			content: []domain.ProductBatch{
				*domain.NewProductBatch(*perishables, 3),
				*domain.NewProductBatch(*equipment, 2),
			},
			want: 3*(10*100) + 2*(120*15) + 200*2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			order, err := domain.NewOrder(tt.dist, tt.transport, tt.content)
			require.NoError(t, err)

			assert.InDelta(t, tt.want, order.CalculateCost(), 1e-9)
		})
	}
}
