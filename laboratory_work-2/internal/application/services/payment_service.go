package services

import (
	"errors"

	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2025/lr-2/internal/domain"
)

var ErrCurrencyMismatch = errors.New("currency mismatch")

type PaymentService struct {
	accounts AccountStorage
}

func NewPaymentService(accounts AccountStorage) *PaymentService {
	return &PaymentService{accounts: accounts}
}

func (s *PaymentService) Deposit(amount domain.Money, accountID uuid.UUID) error {
	account, err := s.accounts.ByID(accountID)
	if err != nil {
		return err
	}

	if account.Balance().Currency() != amount.Currency() {
		return ErrCurrencyMismatch
	}

	if err := account.Deposit(amount); err != nil {
		return err
	}

	return nil
}

func (s *PaymentService) Withdraw(amount domain.Money, accountID uuid.UUID) error {
	account, err := s.accounts.ByID(accountID)
	if err != nil {
		return err
	}

	if account.Balance().Currency() != amount.Currency() {
		return ErrCurrencyMismatch
	}

	if err := account.Withdraw(amount); err != nil {
		return err
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

	if src.Balance().Currency() != t.Value().Currency() || dst.Balance().Currency() != t.Value().Currency() {
		return ErrCurrencyMismatch
	}

	if err := s.Withdraw(t.Value(), t.FromAccountID()); err != nil {
		return err
	}

	if err := s.Deposit(t.Value(), t.ToAccountID()); err != nil {
		if rollbackErr := s.Deposit(t.Value(), t.FromAccountID()); rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}

		return err
	}

	return nil
}
