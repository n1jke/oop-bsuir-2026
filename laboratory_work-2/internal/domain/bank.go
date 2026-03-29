package domain

import "github.com/google/uuid"

type Bank struct {
	id      uuid.UUID
	name    string
	license string
}

func NewBank(id uuid.UUID, name, license string) *Bank {
	return &Bank{id: id, name: name, license: license}
}

func (b Bank) ID() uuid.UUID {
	return b.id
}

func (b Bank) Name() string {
	return b.name
}

func (b Bank) License() string {
	return b.license
}
