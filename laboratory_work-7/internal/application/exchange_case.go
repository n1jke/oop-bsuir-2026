package application

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ExchangeService struct{}

// check: toUserID must exist, ownedBook must exist
func (e *ExchangeService) CreateExchange(ctx context.Context, bookID, toUserID uuid.UUID, expiresAt time.Time, note string) (ExchangeDTO, error) {
	return ExchangeDTO{}, nil
}

// check: exchange must exist, & caller must be either fromUser or toUser.
func (e *ExchangeService) GetExchangeByID(ctx context.Context, exchangeID uuid.UUID) (ExchangeDTO, error) {
	return ExchangeDTO{}, nil
}

// check: exchange must exist and be Pending, caller must be the toUser.
func (e *ExchangeService) AcceptExchange(ctx context.Context, exchangeID uuid.UUID) (ExchangeDTO, error) {
	return ExchangeDTO{}, nil
}

// check: exchange must exist and be Pending, caller must be the toUser.
func (e *ExchangeService) RejectExchange(ctx context.Context, exchangeID uuid.UUID) (ExchangeDTO, error) {
	return ExchangeDTO{}, nil
}

// check: exchange must exist and not be Completed/Canceled, caller must be the fromUser.
func (e *ExchangeService) CancelExchange(ctx context.Context, exchangeID uuid.UUID) (ExchangeDTO, error) {
	return ExchangeDTO{}, nil
}

func (e *ExchangeService) GetUserExchanges(ctx context.Context, userID uuid.UUID, status string) ([]ExchangeDTO, error) {
	return nil, nil
}
