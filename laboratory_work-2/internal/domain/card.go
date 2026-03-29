package domain

import "github.com/google/uuid"

// Card - bank card for account entity.
type Card struct {
	id        uuid.UUID
	number    string
	accountID uuid.UUID
	active    bool
}

func NewCard(id uuid.UUID, number string, accountID uuid.UUID) *Card {
	return &Card{id: id, number: number, accountID: accountID, active: true}
}

func (c Card) ID() uuid.UUID {
	return c.id
}

func (c Card) Number() string {
	return c.number
}

func (c Card) AccountID() uuid.UUID {
	return c.accountID
}

func (c *Card) Activate() {
	c.active = true
}

func (c *Card) Block() {
	c.active = false
}

func (c Card) IsActive() bool {
	return c.active
}
