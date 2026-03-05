package infrastructure

import (
	"fmt"
	"os"
	"time"

	"github.com/n1jke/oop-bsuir-2025/lr-3/internal/domain"
)

// RandomSQLDatabase - имитация тяжелой базы данных.
type RandomSQLDatabase struct {
	connectionString string
}

func NewSQLDatabase(connString ...string) *RandomSQLDatabase {
	if len(connString) == 1 {
		return &RandomSQLDatabase{connectionString: connString[0]}
	}

	return &RandomSQLDatabase{connectionString: "random://root:password@localhost:228/shop"}
}

// Сохранение заказа в "базу данных".
func (db *RandomSQLDatabase) SaveOrder(order *domain.Order, total float64) error {
	fmt.Println("Connecting to RandomSQL at", db.connectionString, "...")
	time.Sleep(500 * time.Millisecond) // Имитация задержки сети

	file, err := os.OpenFile("orders_db.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	record := fmt.Sprintf("[%s] ID: %s | Type: %s | Total: %.2f\n", time.Now().Format(time.RFC3339), order.ID, order.Type, total)
	if _, err := file.WriteString(record); err != nil {
		return err
	}

	fmt.Println("Order saved successfully.")

	return nil
}
