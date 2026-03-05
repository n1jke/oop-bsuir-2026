package infrastructure

import (
	"errors"

	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2025/lr-1/internal/domain"
)

var ErrNotFound = errors.New("not found")

type MemoryAccountStorage struct {
	byID map[uuid.UUID]*domain.Account
}

func NewMemoryAccountStorage() *MemoryAccountStorage {
	return &MemoryAccountStorage{byID: make(map[uuid.UUID]*domain.Account)}
}

func (s *MemoryAccountStorage) Save(account *domain.Account) error {
	s.byID[account.ID()] = account
	return nil
}

func (s *MemoryAccountStorage) ByID(accountID uuid.UUID) (*domain.Account, error) {
	account, exists := s.byID[accountID]
	if !exists {
		return nil, ErrNotFound
	}

	return account, nil
}

func (s *MemoryAccountStorage) UpdateStatus(accountID uuid.UUID, status domain.AccountStatus) error {
	account, exists := s.byID[accountID]
	if !exists {
		return ErrNotFound
	}

	account.ChangeStatus(status)

	return nil
}

type MemoryEventStorage struct {
	events []*domain.Event
}

func NewMemoryEventStorage() *MemoryEventStorage {
	return &MemoryEventStorage{events: make([]*domain.Event, 0, 16)}
}

func (s *MemoryEventStorage) Save(event *domain.Event) error {
	s.events = append(s.events, event)
	return nil
}

func (s *MemoryEventStorage) QueryAll() []*domain.Event {
	return s.events
}
