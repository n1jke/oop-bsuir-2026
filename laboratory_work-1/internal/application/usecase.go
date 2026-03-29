package application

import (
	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2025/lr-1/internal/application/services"
	"github.com/n1jke/oop-bsuir-2025/lr-1/internal/domain"
)

type TransferUseCase struct {
	payments *services.PaymentService
	events   *services.EventService
}

func NewTransferUseCase(payments *services.PaymentService, events *services.EventService) *TransferUseCase {
	return &TransferUseCase{
		payments: payments,
		events:   events,
	}
}

func (uc *TransferUseCase) Execute(t *domain.Transaction) error {
	if err := uc.payments.ProcessTransaction(t); err != nil {
		t.ChangeStatus(domain.Failed)
		_ = uc.events.Publish(domain.NewEvent(uuid.New(), t.ID(), "transaction_failed"))

		return err
	}

	t.ChangeStatus(domain.Completed)

	_ = uc.events.Publish(domain.NewEvent(uuid.New(), t.ID(), "transaction_completed"))

	return nil
}
