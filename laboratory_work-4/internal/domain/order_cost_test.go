package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderCalculateCost(t *testing.T) {
	t.Parallel()

	plane, err := NewTransportInfo("plane", AirTransport, 150, 850)
	require.NoError(t, err)

	tanker, err := NewTransportInfo("tanker", WaterTransport, 2, 35)
	require.NoError(t, err)

	perishables, err := NewCargoInfo("perishables", 10, 100)
	require.NoError(t, err)

	equipment, err := NewCargoInfo("equipment", 120, 15)
	require.NoError(t, err)

	tests := []struct {
		name      string
		dist      float64
		transport Transport
		content   []ProductBatch
		want      float64
	}{
		{
			name:      "only cargo component : distance is zero",
			dist:      0,
			transport: tanker,
			content: []ProductBatch{
				*NewProductBatch(*equipment, 1),
			},
			want: 1 * (120 * 15),
		},
		{
			name:      "cargo with air delivery",
			dist:      250,
			transport: plane,
			content: []ProductBatch{
				*NewProductBatch(*perishables, 2),
				*NewProductBatch(*equipment, 1),
			},
			want: 2*(10*100) + 1*(120*15) + 250*150,
		},
		{
			name:      "cargo with water delivery",
			dist:      200,
			transport: tanker,
			content: []ProductBatch{
				*NewProductBatch(*perishables, 3),
				*NewProductBatch(*equipment, 2),
			},
			want: 3*(10*100) + 2*(120*15) + 200*2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			order, err := NewOrder(tt.dist, tt.transport, tt.content)
			require.NoError(t, err)

			assert.InDelta(t, tt.want, order.CalculateCost(), 1e-9)
		})
	}
}
