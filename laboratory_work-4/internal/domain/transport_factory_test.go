package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCatalogTransportFactory_Create(t *testing.T) {
	truck, err := NewTransportInfo("truck", LandTransport, 15, 80)
	require.NoError(t, err)
	plane, err := NewTransportInfo("plane", AirTransport, 150, 850)
	require.NoError(t, err)

	factory := NewCatalogTransportFactory([]TransportInfo{*truck, *plane})

	tests := []struct {
		name          string
		transportType TransportType
		want          Transport
		wantErr       bool
	}{
		{
			name:          "land valid case",
			transportType: TransportType("truck"),
			want:          &transportModel{*truck},
			wantErr:       false,
		},
		{
			name:          "invalid: not found",
			transportType: TransportType("pumpalumpa"),
			want:          nil,
			wantErr:       true,
		},
		{
			name:          "air valid case",
			transportType: TransportType("plane"),
			want:          &transportModel{*plane},
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := factory.Create(tt.transportType)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
