package infrastructure_test

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-4/internal/domain"
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-4/internal/infrastructure"
)

const (
	validCargo string = `Тип записи;Наименование;Тип;Масса_ед_кг;Стоимость_перевозки_за_кг;Расход_на_км;Скорость_км_ч
cargo;Электроника;;1.5;50;;
cargo;Одежда;;0.8;20;;
transport;Грузовик;land;;;15.0;80`

	invalidCargo = `Тип записи;Наименование;Тип;Масса_ед_кг;Стоимость_перевозки_за_кг;Расход_на_км;Скорость_км_ч
cargo;Электроника;;abc;50;;`

	validTransport = `Тип записи;Наименование;Тип;Масса_ед_кг;Стоимость_перевозки_за_кг;Расход_на_км;Скорость_км_ч
transport;Грузовик;land;;;15.0;80
transport;Самолет;air;;;150.0;850
cargo;Одежда;;0.8;20;;
`

	invalidTransport = `Тип записи;Наименование;Тип;Масса_ед_кг;Стоимость_перевозки_за_кг;Расход_на_км;Скорость_км_ч
transport;Грузовик;unknown;;;15.0;80
`
)

func TestCsvRepository_LoadCargoInfo(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name    string
		path    string
		want    []domain.CargoInfo
		wantErr bool
	}{
		{
			name: "valid",
			path: createTempCSV(t, validCargo),
			want: func() []domain.CargoInfo {
				c1, _ := domain.NewCargoInfo("Электроника", 1.5, 50)
				c2, _ := domain.NewCargoInfo("Одежда", 0.8, 20)

				return []domain.CargoInfo{*c1, *c2}
			}(),
			wantErr: false,
		},
		{
			name:    "invalid",
			path:    createTempCSV(t, invalidCargo),
			want:    []domain.CargoInfo{},
			wantErr: false,
		},
		{
			name:    "file not exists",
			path:    filepath.Join(t.TempDir(), "not_exists.csv"),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := infrastructure.NewCsvRepository(logger, tt.path)
			got, err := c.LoadCargoInfo()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCsvRepository_LoadTransportInfo(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	tests := []struct {
		name    string
		path    string
		want    []domain.TransportInfo
		wantErr bool
	}{
		{
			name: "valid",
			path: createTempCSV(t, validTransport),
			want: func() []domain.TransportInfo {
				t1, _ := domain.NewTransportInfo("Грузовик", domain.LandTransport, 15.0, 80)
				t2, _ := domain.NewTransportInfo("Самолет", domain.AirTransport, 150.0, 850)

				return []domain.TransportInfo{*t1, *t2}
			}(),
			wantErr: false,
		},
		{
			name:    "invalid",
			path:    createTempCSV(t, invalidTransport),
			want:    []domain.TransportInfo{},
			wantErr: false,
		},
		{
			name:    "file not exists",
			path:    filepath.Join(t.TempDir(), "not_exists.csv"),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := infrastructure.NewCsvRepository(logger, tt.path)
			got, err := c.LoadTransportInfo()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func createTempCSV(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "logistic.csv")

	err := os.WriteFile(path, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("failed to write temp csv: %v", err)
	}

	return path
}
