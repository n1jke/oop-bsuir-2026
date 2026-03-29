package domain

import "github.com/google/uuid"

// Account - bank account entity.
type Account struct {
	id       uuid.UUID
	number   string
	clientID uuid.UUID
	currency Currency
	balance  int
	status   AccountStatus
}

// AccountStatus - value object for bank account status.
type AccountStatus string

const (
	AccountActive AccountStatus = "active"
	AccountFrozen AccountStatus = "frozen"
	AccountClosed AccountStatus = "closed"
)

func NewAccount(id uuid.UUID, number string, clientID uuid.UUID, curr Currency) *Account {
	return &Account{
		id:       id,
		number:   number,
		clientID: clientID,
		currency: curr,
		balance:  0,
		status:   AccountActive,
	}
}

func (a *Account) ID() uuid.UUID {
	return a.id
}

func (a *Account) Number() string {
	return a.number
}

func (a *Account) ClientID() uuid.UUID {
	return a.clientID
}

func (a *Account) Balance() Money {
	return NewMoney(a.balance, a.currency)
}

func (a *Account) Status() AccountStatus {
	return a.status
}

func (a *Account) Deposit(amount int) {
	if a.status != AccountActive || amount <= 0 {
		return
	}

	a.balance += amount
}

func (a *Account) Withdraw(amount int) bool {
	if !a.CanWithdraw(amount) {
		return false
	}

	a.balance -= amount

	return true
}

func (a *Account) CanWithdraw(amount int) bool {
	if a.status != AccountActive || amount <= 0 {
		return false
	}

	return a.balance >= amount
}

func (a *Account) IsActive() bool {
	return a.status == AccountActive
}

func (a *Account) ChangeStatus(status AccountStatus) {
	a.status = status
}
