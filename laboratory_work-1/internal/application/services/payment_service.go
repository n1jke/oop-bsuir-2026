package services

import (
	"errors"

	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2025/lr-1/internal/domain"
)

var (
	ErrAccountNotActive = errors.New("account is not active")
	ErrWithdrawRejected = errors.New("withdraw rejected")
	ErrCurrencyMismatch = errors.New("currency mismatch")
)

type PaymentService struct {
	accounts *AccountRepository
}

func NewPaymentService(accounts *AccountRepository) *PaymentService {
	return &PaymentService{accounts: accounts}
}

func (s *PaymentService) Deposit(amount int, accountID uuid.UUID) error {
	account, err := s.accounts.ByID(accountID)
	if err != nil {
		return err
	}

	if !account.IsActive() {
		return ErrAccountNotActive
	}

	account.Deposit(amount)

	return nil
}

func (s *PaymentService) Withdraw(amount int, accountID uuid.UUID) error {
	account, err := s.accounts.ByID(accountID)
	if err != nil {
		return err
	}

	if !account.IsActive() {
		return ErrAccountNotActive
	}

	if !account.Withdraw(amount) {
		return ErrWithdrawRejected
	}

	return nil
}

func (s *PaymentService) ProcessTransaction(t *domain.Transaction) error {
	src, err := s.accounts.ByID(t.FromAccountID())
	if err != nil {
		return err
	}

	dst, err := s.accounts.ByID(t.ToAccountID())
	if err != nil {
		return err
	}

	if src.Balance().Currency() != t.Currency() || dst.Balance().Currency() != t.Currency() {
		return ErrCurrencyMismatch
	}

	if err := s.Withdraw(t.Amount(), t.FromAccountID()); err != nil {
		return err
	}

	if err := s.Deposit(t.Amount(), t.ToAccountID()); err != nil {
		_ = s.Deposit(t.Amount(), t.FromAccountID())
		return err
	}

	return nil
}
