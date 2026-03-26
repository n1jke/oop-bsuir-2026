package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/application"
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/infrastructure"
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/infrastructure/save"
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/infrastructure/stock"
)

func main() {
	logger := slog.Default()
	client := infrastructure.NewCLIOrderSource(logger)

	fileType, filePath, err := client.RequestConfig()
	if err != nil {
		logger.Error("failed to get config input", slog.Any("err", err))
		os.Exit(1)
	}

	stockInfo, err := stock.NewStockLoader(logger, filePath, fileType)
	if err != nil {
		logger.Error("failed to create stock loader", slog.Any("err", err))
		os.Exit(1)
	}

	service, err := application.NewLogisticService(
		application.WithLogger(logger),
		application.WithStock(stockInfo),
		application.WithClient(client),
	)
	if err != nil {
		logger.Error("failed to create logistic service", slog.Any("err", err))
		os.Exit(1)
	}

	resp, err := service.Process()
	if err != nil {
		logger.Error("failed to process order", slog.Any("err", err))
		os.Exit(1)
	}

	sortFields, err := client.RequestSortFields()
	if err != nil {
		logger.Error("failed to get sort fields", slog.Any("err", err))
		os.Exit(1)
	}

	if err := resp.SortFields(sortFields); err != nil {
		logger.Error("failed to sort response options", slog.Any("err", err))
		os.Exit(1)
	}

	exportCfg, err := client.RequestExportConfig()
	if err != nil {
		logger.Error("failed to get export config", slog.Any("err", err))
		os.Exit(1)
	}

	saver, err := save.NewResponseSaver(exportCfg)
	if err != nil {
		logger.Error("failed to create response saver", slog.Any("err", err))
		os.Exit(1)
	}

	if err := saver.Save(resp); err != nil {
		logger.Error("failed to save response", slog.Any("err", err))
		os.Exit(1)
	}

	fmt.Printf("Response saved to: %s\n", exportCfg.OutPath)
}
