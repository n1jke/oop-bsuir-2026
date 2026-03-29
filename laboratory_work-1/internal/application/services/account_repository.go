package services

import (
	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2025/lr-1/internal/domain"
	"github.com/n1jke/oop-bsuir-2025/lr-1/internal/infrastructure"
)

type AccountRepository struct {
	storage *infrastructure.MemoryAccountStorage
}

func NewAccountRepository(storage *infrastructure.MemoryAccountStorage) *AccountRepository {
	return &AccountRepository{storage: storage}
}

func (r *AccountRepository) Create(account *domain.Account) error {
	return r.storage.Save(account)
}

func (r *AccountRepository) ChangeStatus(accountID uuid.UUID, status domain.AccountStatus) error {
	return r.storage.UpdateStatus(accountID, status)
}

func (r *AccountRepository) ByID(accountID uuid.UUID) (*domain.Account, error) {
	return r.storage.ByID(accountID)
}
