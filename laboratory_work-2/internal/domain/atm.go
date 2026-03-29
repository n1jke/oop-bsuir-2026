package domain

import "github.com/google/uuid"

type ATM struct {
	id       uuid.UUID
	branchID uuid.UUID
	cash     int
}

func NewATM(id, branchID uuid.UUID, cash int) *ATM {
	return &ATM{id: id, branchID: branchID, cash: cash}
}

func (a ATM) ID() uuid.UUID {
	return a.id
}

func (a ATM) BranchID() uuid.UUID {
	return a.branchID
}

func (a ATM) Cash() int {
	return a.cash
}

func (a *ATM) Withdraw(amount int) bool {
	if amount <= 0 || a.cash < amount {
		return false
	}

	a.cash -= amount

	return true
}

func (a *ATM) Deposit(amount int) {
	if amount <= 0 {
		return
	}

	a.cash += amount
}
