package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// Employee - bank Employee entity.
type Employee struct {
	id       uuid.UUID
	branchID uuid.UUID
	position string
	employed bool
}

func NewEmployee(branchID uuid.UUID, position string) *Employee {
	return &Employee{
		id:       uuid.New(),
		branchID: branchID,
		position: position,
		employed: true,
	}
}

func (e Employee) ID() uuid.UUID {
	return e.id
}

func (e Employee) BranchID() uuid.UUID {
	return e.branchID
}

func (e Employee) Position() string {
	return e.position
}

func (e *Employee) Hire() {
	e.employed = true
}

func (e *Employee) Fire() {
	e.employed = false
}

func (e Employee) IsEmployed() bool {
	return e.employed
}

func (e Employee) Work() {
	// simulate work
	fmt.Println("do some work")
}
