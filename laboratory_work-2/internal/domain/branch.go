package domain

import "github.com/google/uuid"

// Branch - bank branch entity.
type Branch struct {
	id      uuid.UUID
	bankID  uuid.UUID
	address string
	opened  bool
}

func NewBranch(id, bankID uuid.UUID, address string) *Branch {
	return &Branch{id: id, bankID: bankID, address: address, opened: true}
}

func (b Branch) ID() uuid.UUID {
	return b.id
}

func (b Branch) BankID() uuid.UUID {
	return b.bankID
}

func (b Branch) Address() string {
	return b.address
}

func (b *Branch) Open() {
	b.opened = true
}

func (b *Branch) Close() {
	b.opened = false
}

func (b Branch) IsOpened() bool {
	return b.opened
}
