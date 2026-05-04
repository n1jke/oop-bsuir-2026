package application

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ExchangeService struct{}

func (e *ExchangeService) CreateExchange(ctx context.Context, bookID, toUserID uuid.UUID, expiresAt time.Time, note string) (ExchangeDTO, error) {
	return ExchangeDTO{}, nil
}

func (e *ExchangeService) GetExchangeByID(ctx context.Context, exchangeID uuid.UUID) (ExchangeDTO, error) {
	return ExchangeDTO{}, nil
}

func (e *ExchangeService) AcceptExchange(ctx context.Context, exchangeID uuid.UUID) (ExchangeDTO, error) {
	return ExchangeDTO{}, nil
}

func (e *ExchangeService) RejectExchange(ctx context.Context, exchangeID uuid.UUID) (ExchangeDTO, error) {
	return ExchangeDTO{}, nil
}

func (e *ExchangeService) CancelExchange(ctx context.Context, exchangeID uuid.UUID) (ExchangeDTO, error) {
	return ExchangeDTO{}, nil
}

func (e *ExchangeService) GetUserExchanges(ctx context.Context, userID uuid.UUID, status string) ([]ExchangeDTO, error) {
	return nil, nil
}
