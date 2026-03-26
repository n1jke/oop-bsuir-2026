package infrastructure

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/application"
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/domain"
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/infrastructure/save"
)

var (
	ErrEmptyCargoCatalog     = errors.New("cargo catalog is empty")
	ErrEmptyTransportCatalog = errors.New("transport catalog is empty")
	ErrInvalidMenuSelection  = errors.New("invalid menu selection")
	ErrInvalidDistanceInput  = errors.New("invalid distance input")
	ErrInvalidItemsCount     = errors.New("invalid items count")
	ErrInvalidBatchCount     = errors.New("invalid batch count")
	ErrInvalidFileTypeInput  = errors.New("invalid file type input")
	ErrEmptyFilePathInput    = errors.New("empty file path input")
	ErrInvalidExportFormat   = errors.New("invalid export format input")
	ErrInvalidYesNoInput     = errors.New("invalid yes/no input")
)

type CLIOrderSource struct {
	logger *slog.Logger
}

func NewCLIOrderSource(logger *slog.Logger) *CLIOrderSource {
	return &CLIOrderSource{
		logger: logger,
	}
}

func (c *CLIOrderSource) RequestConfig() (string, string, error) {
	raw, err := ReadLine("Enter config file type (csv/json/xml)")
	if err != nil {
		return "", "", err
	}

	fileType := strings.ToLower(strings.TrimSpace(raw))
	if fileType != "csv" && fileType != "json" && fileType != "xml" {
		return "", "", ErrInvalidFileTypeInput
	}

	rawPath, err := ReadLine("Enter config file path")
	if err != nil {
		return "", "", err
	}

	filePath := strings.TrimSpace(rawPath)
	if filePath == "" {
		return "", "", ErrEmptyFilePathInput
	}

	return fileType, filePath, nil
}

func (c *CLIOrderSource) RequestExportConfig() (save.ExportConfig, error) {
	rawFormat, err := ReadLine("Enter export format (json/yaml): ")
	if err != nil {
		return save.ExportConfig{}, err
	}

	format := strings.ToLower(strings.TrimSpace(rawFormat))
	if format != "json" && format != "yaml" {
		return save.ExportConfig{}, ErrInvalidExportFormat
	}

	useEncryption, err := readYesNo("Enable encryption (y/n): ")
	if err != nil {
		return save.ExportConfig{}, err
	}

	useCompression, err := readYesNo("Enable compression (y/n): ")
	if err != nil {
		return save.ExportConfig{}, err
	}

	transformations := make([]string, 0, 2)
	if useEncryption {
		transformations = append(transformations, "encrypt")
	}

	if useCompression {
		transformations = append(transformations, "compress")
	}

	outPath := "response." + format
	if useEncryption {
		outPath += ".enc"
	}

	if useCompression {
		outPath += ".gz"
	}

	return save.ExportConfig{
		Format:          format,
		Transformations: transformations,
		OutPath:         outPath,
	}, nil
}

func (c *CLIOrderSource) RequestOrder(cargo []domain.CargoInfo, transport []domain.TransportInfo) (*application.ClientResponse, error) {
	printTransportCatalog(transport)

	transportIdx, err := ReadInt("Select transport (or -1 for all): ")
	if err != nil {
		return nil, err
	}

	if transportIdx < -1 || transportIdx >= len(transport) {
		return nil, ErrInvalidMenuSelection
	}

	factory := domain.NewCatalogTransportFactory(transport)

	var selectedTransport domain.Transport
	if transportIdx != -1 {
		selectedTransport, err = factory.Create(transport[transportIdx].Name())
		if err != nil {
			return nil, err
		}
	}

	distance, err := ReadFloat("Input distance: ")
	if err != nil {
		return nil, err
	}

	if distance < 0 {
		return nil, ErrInvalidDistanceInput
	}

	printCargoCatalog(cargo)

	itemsCount, err := ReadInt("Input count of items: ")
	if err != nil {
		return nil, err
	}

	if itemsCount <= 0 {
		return nil, ErrInvalidItemsCount
	}

	content := make([]domain.ProductBatch, 0, itemsCount)
	for i := range itemsCount {
		cargoIdx, err := ReadInt(fmt.Sprintf("Pos %d, select item: ", i+1))
		if err != nil {
			return nil, err
		}

		if cargoIdx < 0 || cargoIdx >= len(cargo) {
			return nil, ErrInvalidMenuSelection
		}

		count, err := ReadInt(fmt.Sprintf("Pos %d, select count: ", i+1))
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

func ReadLine(msg string) (string, error) {
	fmt.Print(msg)

	line := ""
	_, _ = fmt.Scan(&line)

	return strings.TrimSpace(line), nil
}

func ReadInt(msg string) (int, error) {
	raw, err := ReadLine(msg)
	if err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", ErrInvalidMenuSelection, raw)
	}

	return value, nil
}

func ReadFloat(msg string) (float64, error) {
	raw, err := ReadLine(msg)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", ErrInvalidDistanceInput, raw)
	}

	return value, nil
}

func readYesNo(msg string) (bool, error) {
	raw, err := ReadLine(msg)
	if err != nil {
		return false, err
	}

	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return false, ErrInvalidYesNoInput
	}
}
