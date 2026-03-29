package application_test

//go:generate go run go.uber.org/mock/mockgen@latest -source=order_logistic.go -destination=mocks/mock.go -package=mocks StockRepository,OrderRequestSource

import (
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-4/internal/application"
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-4/internal/application/mocks"
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-4/internal/domain"
)

func TestLogisticServiceProcess_TableDriven(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	dist := 100.0
	tempErr := errors.New("temp error")

	// setup
	cargoCatalog := cargoCatalog(t)
	transportCatalog := transportCatalog(t)
	factory := domain.NewCatalogTransportFactory(transportCatalog)
	transport, _ := factory.Create(transportCatalog[0].Name())
	content := []domain.ProductBatch{
		*domain.NewProductBatch(cargoCatalog[0], 10),
	}

	tests := []struct {
		name          string
		prepareMocks  func(store *mocks.MockStockRepository, client *mocks.MockOrderRequestSource)
		wantErr       bool
		checkResponse func(t *testing.T, got *application.ServiceResponse)
	}{
		{
			name: "cargo load error",
			prepareMocks: func(store *mocks.MockStockRepository, client *mocks.MockOrderRequestSource) {
				store.EXPECT().LoadCargoInfo().Return(nil, tempErr).Times(1)
				store.EXPECT().LoadTransportInfo().Times(0)
				client.EXPECT().RequestOrder(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
		},
		{
			name: "transport load error",
			prepareMocks: func(store *mocks.MockStockRepository, client *mocks.MockOrderRequestSource) {
				store.EXPECT().LoadCargoInfo().Return(cargoCatalog, nil).Times(1)
				store.EXPECT().LoadTransportInfo().Return(nil, tempErr).Times(1)
				client.EXPECT().RequestOrder(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
		},
		{
			name: "request order error",
			prepareMocks: func(store *mocks.MockStockRepository, client *mocks.MockOrderRequestSource) {
				store.EXPECT().LoadCargoInfo().Return(cargoCatalog, nil).Times(1)
				store.EXPECT().LoadTransportInfo().Return(transportCatalog, nil).Times(1)
				client.EXPECT().RequestOrder(gomock.Any(), gomock.Any()).
					Return(nil, tempErr).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "happy end",
			prepareMocks: func(store *mocks.MockStockRepository, client *mocks.MockOrderRequestSource) {
				store.EXPECT().LoadCargoInfo().Return(cargoCatalog, nil).Times(1)
				store.EXPECT().LoadTransportInfo().Return(transportCatalog, nil).Times(1)
				client.EXPECT().RequestOrder(cargoCatalog, transportCatalog).
					Return(&application.ClientResponse{
						Transport: transport,
						Dist:      dist,
						Content:   content,
					}, nil).
					Times(1)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// give
			store := mocks.NewMockStockRepository(ctrl)
			client := mocks.NewMockOrderRequestSource(ctrl)
			tt.prepareMocks(store, client)

			srv, err := application.NewLogisticService(
				application.WithLogger(logger),
				application.WithStore(store),
				application.WithClient(client),
			)
			require.NoError(t, err)

			// when
			got, err := srv.Process()

			// then
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}

func cargoCatalog(t *testing.T) []domain.CargoInfo {
	t.Helper()

	electronics, err := domain.NewCargoInfo("Электроника", 1.5, 50)
	require.NoError(t, err)
	clothes, err := domain.NewCargoInfo("Одежда", 0.8, 20)
	require.NoError(t, err)

	return []domain.CargoInfo{*electronics, *clothes}
}

func transportCatalog(t *testing.T) []domain.TransportInfo {
	t.Helper()

	truck, err := domain.NewTransportInfo("Грузовик", domain.LandTransport, 15, 80)
	require.NoError(t, err)
	plane, err := domain.NewTransportInfo("Самолет", domain.AirTransport, 150, 850)
	require.NoError(t, err)

	return []domain.TransportInfo{*truck, *plane}
}
