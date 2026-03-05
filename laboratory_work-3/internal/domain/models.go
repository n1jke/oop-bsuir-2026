package domain

import (
	"errors"
)

var (
	ErrItemsCount       = errors.New("order must have at least one item")
	ErrRequireCity      = errors.New("destination city is required")
	ErrUnknowOrderType  = errors.New("unknown order type")
	ErrBudgetOrderSizze = errors.New("budget orders cannot have more than 3 items")
)

// Item - товар в заказе.
type Item struct {
	ID    string
	Name  string
	Price float64
}

// Address - адрес доставки.
type Address struct {
	City   string
	Street string
	Zip    string
}

// Order - заказ.
type Order struct {
	ID              string
	Items           []Item
	Type            OrderType
	DiscountProgram Discount
	ClientEmail     string
	Destination     Address
}

func (o *Order) Validate() error {
	if len(o.Items) == 0 {
		return ErrItemsCount
	}

	if o.Destination.City == "" {
		return ErrRequireCity
	}

	return nil
}

type Discount string

const (
	Gold   Discount = "Gold"
	Silver Discount = "Silver"
	Newbie Discount = "Newbie"
)

func (d Discount) CalculateDiscount(base float64) float64 {
	switch d {
	case "Gold":
		return base * 0.85
	case "Silver":
		return base * 0.90
	default:
		return base
	}
}

type OrderType string

const (
	Standart      OrderType = "Standart"
	Premium       OrderType = "Premium"
	Budget        OrderType = "Budget"
	International OrderType = "International"
)

func (ot OrderType) CalculatePrice(base float64, items []Item) (float64, error) {
	switch ot {
	case Standart:
		base *= 1.2
	case Premium:
		base = (base * 0.9) * 1.2
	case Budget:
		if len(items) > 3 {
			return 0, ErrBudgetOrderSizze
		}
	case International:
		base *= 1.5
	default:
		return 0, ErrUnknowOrderType
	}

	return base, nil
}
