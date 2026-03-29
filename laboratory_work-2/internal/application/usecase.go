package application

import (
	"errors"

	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2025/lr-2/internal/application/services"
	"github.com/n1jke/oop-bsuir-2025/lr-2/internal/domain"
)

var (
	ErrNilTransaction     = errors.New("transaction is nil")
	ErrInvalidBonusPoints = errors.New("bonus points must be positive")
	ErrNotSavingsAccount  = errors.New("bonus redeem is allowed only for savings accounts")
)

type TransferUseCase struct {
	payments *services.PaymentService
	events   services.EventStorage
}

func NewTransferUseCase(
	payments *services.PaymentService,
	events services.EventStorage,
) *TransferUseCase {
	return &TransferUseCase{
		payments: payments,
		events:   events,
	}
}

func (uc *TransferUseCase) Execute(t *domain.Transaction) error {
	if t == nil {
		return ErrNilTransaction
	}

	t.ChangeStatus(domain.Pending)

	if err := uc.payments.ProcessTransaction(t); err != nil {
		t.ChangeStatus(domain.Failed)

		failEvent := domain.NewEvent(uuid.New(), t.ID(), "transaction_failed")
		if publishErr := uc.events.Save(failEvent.ID(), *failEvent); publishErr != nil {
			return errors.Join(err, publishErr)
		}

		return err
	}

	t.ChangeStatus(domain.Completed)

	doneEvent := domain.NewEvent(uuid.New(), t.ID(), "transaction_completed")
	if err := uc.events.Save(doneEvent.ID(), *doneEvent); err != nil {
		return err
	}

	return nil
}

type BonusRedeemCommand struct {
	AccountID uuid.UUID
	Points    int
}

type savingsBonusAccount interface {
	services.PaymentAccount
	RedeemBonus(points int) error
}

type BonusRedeemUseCase struct {
	accounts services.AccountStorage
	events   services.EventStorage
}

func NewBonusRedeemUseCase(accounts services.AccountStorage, events services.EventStorage) *BonusRedeemUseCase {
	return &BonusRedeemUseCase{
		accounts: accounts,
		events:   events,
	}
}

func (uc *BonusRedeemUseCase) Execute(command BonusRedeemCommand) error {
	if command.Points <= 0 {
		return ErrInvalidBonusPoints
	}

	account, err := uc.accounts.ByID(command.AccountID)
	if err != nil {
		return err
	}

	savings, ok := account.(savingsBonusAccount)
	if !ok {
		failEvent := domain.NewEvent(uuid.New(), account.ID(), "bonus_redeem_failed")
		if publishErr := uc.events.Save(failEvent.ID(), *failEvent); publishErr != nil {
			return errors.Join(ErrNotSavingsAccount, publishErr)
		}

		return ErrNotSavingsAccount
	}

	if err := savings.RedeemBonus(command.Points); err != nil {
		failEvent := domain.NewEvent(uuid.New(), account.ID(), "bonus_redeem_failed")
		if publishErr := uc.events.Save(failEvent.ID(), *failEvent); publishErr != nil {
			return errors.Join(err, publishErr)
		}

		return err
	}

	doneEvent := domain.NewEvent(uuid.New(), account.ID(), "bonus_redeemed")

	return uc.events.Save(doneEvent.ID(), *doneEvent)
}
