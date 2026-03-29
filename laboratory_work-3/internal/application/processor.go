package application

import (
	"fmt"

	"github.com/n1jke/oop-bsuir-2025/lr-3/internal/domain"
)

type Repository interface {
	SaveOrder(order *domain.Order, total float64) error
}

type CacheService interface {
	TryAddOrder(order *domain.Order) bool
	RemoveOrder(id string)
	FindOrder(id string) (*domain.Order, bool)
}

type Notifier interface {
	Notify(to string, subject string, body string)
}

type OrderProcessor struct {
	database   Repository
	cache      CacheService
	clientMsg  Notifier
	managerMsg Notifier
}

func NewOrderProcessor(db Repository, cache CacheService, cMsg, mMsg Notifier) *OrderProcessor {
	return &OrderProcessor{
		database:   db,
		cache:      cache,
		clientMsg:  cMsg,
		managerMsg: mMsg,
	}
}

func (op *OrderProcessor) Process(order *domain.Order) error {
	fmt.Printf("--- Processing Order %s ---\n", order.ID)

	err := order.Validate()
	if err != nil {
		return fmt.Errorf("invalid order: %w", err)
	}

	total, err := calculateTotal(order)
	if err != nil {
		return fmt.Errorf("invalid order: %w", err)
	}

	err = op.processCacheAndSave(order, total)
	if err != nil {
		return fmt.Errorf("failed to save order after processing: %w", err)
	}

	op.sendNotifications(order, total)

	return nil
}

func calculateTotal(order *domain.Order) (float64, error) {
	var base float64
	for _, item := range order.Items {
		base += item.Price
	}

	base = order.DiscountProgram.CalculateDiscount(base)

	total, err := order.Type.CalculatePrice(base, order.Items)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (op *OrderProcessor) processCacheAndSave(order *domain.Order, total float64) error {
	if !op.cache.TryAddOrder(order) {
		fmt.Println("Order is already processed")
		return nil
	}

	if err := op.database.SaveOrder(order, total); err != nil {
		op.cache.RemoveOrder(order.ID)
		return fmt.Errorf("database error: %w", err)
	}

	return nil
}

func (op *OrderProcessor) sendNotifications(order *domain.Order, total float64) {
	emailBody := fmt.Sprintf("<h1>Your order %s is confirmed!</h1><p>Total: %.2f</p>", order.ID, total)
	telegramBody := fmt.Sprintf("<h1>Order for client %s with orderId %s is confirmed!</h1><p>Total: %.2f</p>", order.ClientEmail,
		order.ID, total)

	op.clientMsg.Notify(order.ClientEmail, "Order Confirmation", emailBody)
	op.managerMsg.Notify("manager", "Order Notification", telegramBody)
}
