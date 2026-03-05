package infrastructure

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-4/internal/application"
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-4/internal/domain"
)

var (
	ErrEmptyCargoCatalog     = errors.New("cargo catalog is empty")
	ErrEmptyTransportCatalog = errors.New("transport catalog is empty")
	ErrInvalidMenuSelection  = errors.New("invalid menu selection")
	ErrInvalidDistanceInput  = errors.New("invalid distance input")
	ErrInvalidItemsCount     = errors.New("invalid items count")
	ErrInvalidBatchCount     = errors.New("invalid batch count")
)

type CLIOrderSource struct {
	logger *slog.Logger
}

func NewCLIOrderSource(logger *slog.Logger) *CLIOrderSource {
	return &CLIOrderSource{
		logger: logger,
	}
}

func (c *CLIOrderSource) RequestOrder(cargo []domain.CargoInfo, transport []domain.TransportInfo) (*application.ClientResponse, error) {
	if len(cargo) == 0 {
		return nil, ErrEmptyCargoCatalog
	}

	if len(transport) == 0 {
		return nil, ErrEmptyTransportCatalog
	}

	printTransportCatalog(transport)

	transportIdx, err := c.readInt("Select transport: ")
	if err != nil {
		return nil, err
	}

	if transportIdx < 0 || transportIdx >= len(transport) {
		return nil, ErrInvalidMenuSelection
	}

	factory := domain.NewCatalogTransportFactory(transport)

	selectedTransport, err := factory.Create(transport[transportIdx].Name())
	if err != nil {
		return nil, err
	}

	distance, err := c.readFloat("Input distance: ")
	if err != nil {
		return nil, err
	}

	if distance < 0 {
		return nil, ErrInvalidDistanceInput
	}

	printCargoCatalog(cargo)

	itemsCount, err := c.readInt("Input count of items: ")
	if err != nil {
		return nil, err
	}

	if itemsCount <= 0 {
		return nil, ErrInvalidItemsCount
	}

	content := make([]domain.ProductBatch, 0, itemsCount)
	for i := range itemsCount {
		cargoIdx, err := c.readInt(fmt.Sprintf("Pos %d, select item: ", i+1))
		if err != nil {
			return nil, err
		}

		if cargoIdx < 0 || cargoIdx >= len(cargo) {
			return nil, ErrInvalidMenuSelection
		}

		count, err := c.readInt(fmt.Sprintf("Pos %d, select count: ", i+1))
		if err != nil {
			return nil, err
		}

		if count <= 0 {
			return nil, ErrInvalidBatchCount
		}

		content = append(content, *domain.NewProductBatch(cargo[cargoIdx], uint(count)))
	}

	return &application.ClientResponse{
		Transport: selectedTransport,
		Dist:      distance,
		Content:   content,
	}, nil
}

func printTransportCatalog(transport []domain.TransportInfo) {
	fmt.Println("Available transport:")

	for i := range transport {
		fmt.Printf(
			"[%d] %s | mode=%s | rate=%.2f/km| speed=%.2f km/h\n",
			i,
			transport[i].Name(),
			transport[i].Mode(),
			transport[i].DeliveryRate(),
			transport[i].Speed(),
		)
	}
}

func printCargoCatalog(cargo []domain.CargoInfo) {
	fmt.Println("Available cargo:")

	for i := range cargo {
		fmt.Printf(
			"[%d] %s | unitWeight=%.2f kg | costPerKg=%.2f\n",
			i,
			cargo[i].Name(),
			cargo[i].Weight(),
			cargo[i].CostPerKg(),
		)
	}
}

func (c *CLIOrderSource) readLine(msg string) (string, error) {
	fmt.Print(msg + ":")

	line := ""
	_, _ = fmt.Scan(&line)

	return strings.TrimSpace(line), nil
}

func (c *CLIOrderSource) readInt(msg string) (int, error) {
	raw, err := c.readLine(msg)
	if err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", ErrInvalidMenuSelection, raw)
	}

	return value, nil
}

func (c *CLIOrderSource) readFloat(msg string) (float64, error) {
	raw, err := c.readLine(msg)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", ErrInvalidDistanceInput, raw)
	}

	return value, nil
}
