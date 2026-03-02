package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-4/internal/application"
	"github.com/n1jke/oop-bsuir-2025/laboratory_work-4/internal/infrastructure"
)

const logisticFile string = "laboratory_work-4/logistic.csv"

func main() {
	logger := slog.Default()
	store := infrastructure.NewCsvRepository(logger, logisticFile)
	client := infrastructure.NewCLIOrderSource(logger)

	service, err := application.NewLogisticService(
		application.WithLogger(logger),
		application.WithStore(store),
		application.WithClient(client),
	)
	if err != nil {
		logger.Error("failed to create logistic service", "error", err)
		os.Exit(1)
	}

	resp, err := service.Process()
	if err != nil {
		logger.Error("failed to process order", "error", err)
		os.Exit(1)
	}

	fmt.Printf("ID: %s\n", resp.Order.ID())
	fmt.Printf("Total cost: %.2f\n", resp.Cost)
	fmt.Printf("Delivery duration: %s\n", resp.Order.EstimateDuration())
}
