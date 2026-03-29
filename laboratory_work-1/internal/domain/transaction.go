package domain

import "github.com/google/uuid"

// Transaction - transfer between two accounts.
type Transaction struct {
	id            uuid.UUID
	fromAccountID uuid.UUID
	toAccountID   uuid.UUID
	amount        int
	currency      Currency
	status        TransactionStatus
}

// TransactionStatus - value object for transaction status.
type TransactionStatus string

const (
	Initiated TransactionStatus = "initiated"
	Completed TransactionStatus = "completed"
	Failed    TransactionStatus = "failed"
)

func NewTransaction(id, fromAccountID, toAccountID uuid.UUID, amount int, curr Currency) *Transaction {
	return &Transaction{
		id:            id,
		fromAccountID: fromAccountID,
		toAccountID:   toAccountID,
		amount:        amount,
		currency:      curr,
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

func (t *Transaction) Amount() int {
	return t.amount
}

func (t *Transaction) Currency() Currency {
	return t.currency
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
