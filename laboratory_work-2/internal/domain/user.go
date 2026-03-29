package domain

import "github.com/google/uuid"

// Client - bank customer entity.
type Client struct {
	id       uuid.UUID
	passport string
	fullName string
	active   bool
}

func NewClient(id uuid.UUID, passport, fullName string) *Client {
	return &Client{
		id:       id,
		passport: passport,
		fullName: fullName,
		active:   true,
	}
}

func (c Client) ID() uuid.UUID {
	return c.id
}

func (c Client) Passport() string {
	return c.passport
}

func (c Client) FullName() string {
	return c.fullName
}

func (c *Client) Activate() {
	c.active = true
}

func (c *Client) Deactivate() {
	c.active = false
}

func (c Client) IsActive() bool {
	return c.active
}
