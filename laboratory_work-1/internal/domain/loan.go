package domain

import "github.com/google/uuid"

// Loan - loan deal for client.
type Loan struct {
	id       uuid.UUID
	clientID uuid.UUID
	amount   int
	approved bool
}

func NewLoan(id, clientID uuid.UUID, amount int) *Loan {
	return &Loan{id: id, clientID: clientID, amount: amount, approved: false}
}

func (l Loan) ID() uuid.UUID {
	return l.id
}

func (l Loan) ClientID() uuid.UUID {
	return l.clientID
}

func (l Loan) Amount() int {
	return l.amount
}

func (l *Loan) Approve() {
	l.approved = true
}

func (l *Loan) Reject() {
	l.approved = false
}

func (l Loan) IsApproved() bool {
	return l.approved
}
