package domain

import "fmt"

// Money - {amount of money and currency}.
type Money struct {
	amount int
	curr   Currency
}

// Currency - value object for currency.
type Currency string

func NewMoney(amount int, curr Currency) Money {
	return Money{amount: amount, curr: curr}
}

func (m Money) Amount() int {
	return m.amount
}

func (m Money) Currency() Currency {
	return m.curr
}

func (m *Money) Add(other Money) {
	// todo: handle differenc currency
	m.amount += other.amount
}

func (m *Money) Sub(other Money) {
	// todo: handle differenc currency
	m.amount -= other.amount
}

func (m Money) String() string {
	return fmt.Sprintf("%d %s", m.amount, m.curr)
}
