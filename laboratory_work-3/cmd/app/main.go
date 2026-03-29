package main

import (
	"fmt"
	"log"

	"github.com/n1jke/oop-bsuir-2025/lr-3/internal/application"
	"github.com/n1jke/oop-bsuir-2025/lr-3/internal/domain"
	"github.com/n1jke/oop-bsuir-2025/lr-3/internal/infrastructure"
)

func main() {
	// 0. Инфраструктура
	db := infrastructure.NewSQLDatabase()
	cache := application.NewLocalCacheService()
	clienMsg := infrastructure.NewSMTPMailer("smtp.google.com")
	managerMsg := infrastructure.NewTelegramMailer("adifdhdf")

	// 1. Создание заказа
	order := &domain.Order{
		ID:              "ORD-256-X",
		Type:            domain.Premium,
		DiscountProgram: domain.Gold,
		Items: []domain.Item{
			{ID: "1", Name: "Thermal Clips", Price: 1500},
			{ID: "2", Name: "UNATCO Pass Card", Price: 50},
		},
		ClientEmail: "jeevacation@gmail.com",
		Destination: domain.Address{City: "Agartha", Street: "33 Thomas Street", Zip: "[REDACTED]"},
	}

	// 2. Инициализация процессора
	processor := application.NewOrderProcessor(db, cache, clienMsg, managerMsg)

	// 3. Обработка заказа
	if err := processor.Process(order); err != nil {
		log.Fatalf("Failed to process order: %v", err)
	}
	// checck cache for order in current session
	if err := processor.Process(order); err != nil {
		log.Fatalf("Failed to process order: %v", err)
	}

	// 4. Работа с обслуживанием
	fmt.Println("\nTesting Warehouse Stuff:")

	workers := []application.WarehouseWorker{
		domain.HumanManager{},
		domain.RobotPacker{Model: "George Droid"},
	}

	application.ManageWarehouse(workers)
}
