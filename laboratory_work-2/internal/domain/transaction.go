package domain

import "github.com/google/uuid"

// Transaction - transfer between two accounts.
type Transaction struct {
	id            uuid.UUID
	fromAccountID uuid.UUID
	toAccountID   uuid.UUID
	value         Money
	status        TransactionStatus
}

// TransactionStatus - value object for transaction status.
type TransactionStatus int

const (
	Initiated TransactionStatus = iota
	Pending
	Declined
	Failed
	Completed
)

func NewTransaction(id, fromAccountID, toAccountID uuid.UUID, value Money) *Transaction {
	return &Transaction{
		id:            id,
		fromAccountID: fromAccountID,
		toAccountID:   toAccountID,
		value:         value,
		status:        Initiated,
	}
}

func (t *Transaction) ID() uuid.UUID {
	return t.id
}

func (t *Transaction) FromAccountID() uuid.UUID {
	return t.fromAccountID
}

func (t *Transaction) ToAccountID() uuid.UUID {
	return t.toAccountID
}

func (t *Transaction) Value() Money {
	return t.value
}

func (t *Transaction) Status() TransactionStatus {
	return t.status
}

func (t *Transaction) ChangeStatus(status TransactionStatus) {
	t.status = status
}

func (t *Transaction) IsCompleted() bool {
	return t.status == Completed
}
