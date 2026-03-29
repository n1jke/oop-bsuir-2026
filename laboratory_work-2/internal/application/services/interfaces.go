package services

import (
	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2025/lr-2/internal/domain"
)

type PaymentAccount interface {
	ID() uuid.UUID
	ClientID() uuid.UUID
	Status() domain.AccountStatus
	Balance() domain.Money

	Deposit(m domain.Money) error
	Withdraw(m domain.Money) error
	CanWithdraw(m domain.Money) error
	ChangeStatus(s domain.AccountStatus)
}

type Storage[K comparable, V any] interface {
	Save(k K, v V) error
	ByID(k K) (V, error)
	Delete(k K) error
}

type AccountStorage interface {
	Storage[uuid.UUID, PaymentAccount]
	UpdateStatus(accountUUID uuid.UUID, status domain.AccountStatus) error
}

type EventStorage interface {
	Storage[uuid.UUID, domain.Event]
	QueryAll() []*domain.Event
}
