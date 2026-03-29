package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/domain"
)

func TestCatalogTransportFactory_Create(t *testing.T) {
	truck, err := domain.NewTransportInfo("truck", domain.LandTransport, 15, 80)
	require.NoError(t, err)
	plane, err := domain.NewTransportInfo("plane", domain.AirTransport, 150, 850)
	require.NoError(t, err)

	factory := domain.NewCatalogTransportFactory([]domain.TransportInfo{*truck, *plane})

	tests := []struct {
		name          string
		transportType domain.TransportType
		wantType      domain.TransportType
		wantMode      domain.TransportMode
		wantErr       bool
	}{
		{
			name:          "land valid case",
			transportType: domain.TransportType("truck"),
			wantType:      truck.Name(),
			wantMode:      truck.Mode(),
			wantErr:       false,
		},
		{
			name:          "invalid: not found",
			transportType: domain.TransportType("pumpalumpa"),
			wantErr:       true,
		},
		{
			name:          "air valid case",
			transportType: domain.TransportType("plane"),
			wantType:      plane.Name(),
			wantMode:      plane.Mode(),
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := factory.Create(tt.transportType)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantType, got.Name())
			assert.Equal(t, tt.wantMode, got.Mode())
		})
	}
}
