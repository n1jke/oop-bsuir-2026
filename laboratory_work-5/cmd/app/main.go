package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/application"
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/infrastructure"
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

	fmt.Printf("ID: %s\n", resp.Order.ID())
	fmt.Printf("Total cost: %.2f\n", resp.Cost)
	fmt.Printf("Delivery duration: %s\n", resp.Order.EstimateDuration())
}
