package application

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/domain"
)

type UserIDKeyType struct{}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(UserIDKeyType{}).(uuid.UUID)
	return id, ok
}

type ExchangeService struct {
	logger        *slog.Logger
	exchangeRepo  ExchangeRepository
	ownedBookRepo OwnedBookRepository
	userRepo      UserRepository
	tx            Transactor
}

func NewExchangeService(logger *slog.Logger, exchangeRepo ExchangeRepository, ownedBookRepo OwnedBookRepository,
	userRepo UserRepository, tx Transactor,
) *ExchangeService {
	return &ExchangeService{
		logger:        logger.With(slog.String("module", "exchange-service")),
		exchangeRepo:  exchangeRepo,
		ownedBookRepo: ownedBookRepo,
		userRepo:      userRepo,
		tx:            tx,
	}
}

func (e *ExchangeService) CreateExchange(ctx context.Context, ownedBookID, toUserID uuid.UUID, expiresAt time.Time,
	note string,
) (ExchangeDTO, error) {
	fromID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return ExchangeDTO{}, ErrInvalidCredentials
	}

	var exDTO ExchangeDTO

	err := e.tx.WithTransaction(ctx, func(ctx context.Context) error {
		if _, err := e.userRepo.GetByID(ctx, toUserID); err != nil {
			if errors.Is(err, ErrNotFound) {
				return ErrUserNotFound
			}

			e.logger.Error("get to user", slog.Any("err", err))

			return ErrUnvailible
		}

		now := time.Now()
		dp := domain.NewDatePolicy(now, now, expiresAt)

		exReq, err := domain.NewExchangeRequest(ownedBookID, fromID, toUserID, dp, note)
		if err != nil {
			return ErrInvalidParams
		}

		if err := e.exchangeRepo.Add(ctx, exReq); err != nil {
			return err
		}

		exDTO = mapDomainExchangeToDTO(exReq)

		return nil
	})
	if err != nil {
		return ExchangeDTO{}, err
	}

	return exDTO, nil
}

func (e *ExchangeService) GetExchangeByID(ctx context.Context, exchangeID uuid.UUID) (ExchangeDTO, error) {
	caller, ok := GetUserIDFromContext(ctx)
	if !ok {
		return ExchangeDTO{}, ErrInvalidCredentials
	}

	ex, err := e.exchangeRepo.GetByID(ctx, exchangeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ExchangeDTO{}, ErrNotFound
		}

		e.logger.Error("get exchange", slog.Any("err", err))

		return ExchangeDTO{}, ErrUnvailible
	}

	if ex.FromID() != caller && ex.ToID() != caller {
		return ExchangeDTO{}, ErrInvalidCredentials
	}

	return mapDomainExchangeToDTO(&ex), nil
}

func (e *ExchangeService) AcceptExchange(ctx context.Context, exchangeID uuid.UUID) (ExchangeDTO, error) {
	caller, ok := GetUserIDFromContext(ctx)
	if !ok {
		return ExchangeDTO{}, ErrInvalidCredentials
	}

	var result ExchangeDTO

	err := e.tx.WithTransaction(ctx, func(ctx context.Context) error {
		ex, err := e.exchangeRepo.GetByID(ctx, exchangeID)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return ErrNotFound
			}

			return ErrUnvailible
		}

		if ex.Status() != domain.Pending {
			return ErrInvalidParams
		}

		if ex.ToID() != caller {
			return ErrInvalidCredentials
		}

		updated, err := e.exchangeRepo.UpdateStatus(ctx, exchangeID, domain.Accepted)
		if err != nil {
			return err
		}

		result = mapDomainExchangeToDTO(&updated)

		return nil
	})
	if err != nil {
		return ExchangeDTO{}, err
	}

	return result, nil
}

func (e *ExchangeService) RejectExchange(ctx context.Context, exchangeID uuid.UUID) (ExchangeDTO, error) {
	caller, ok := GetUserIDFromContext(ctx)
	if !ok {
		return ExchangeDTO{}, ErrInvalidCredentials
	}

	var result ExchangeDTO

	err := e.tx.WithTransaction(ctx, func(ctx context.Context) error {
		ex, err := e.exchangeRepo.GetByID(ctx, exchangeID)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return ErrNotFound
			}

			return ErrUnvailible
		}

		if ex.Status() != domain.Pending {
			return ErrInvalidParams
		}

		if ex.ToID() != caller {
			return ErrInvalidCredentials
		}

		updated, err := e.exchangeRepo.UpdateStatus(ctx, exchangeID, domain.Rejected)
		if err != nil {
			return err
		}

		result = mapDomainExchangeToDTO(&updated)

		return nil
	})
	if err != nil {
		return ExchangeDTO{}, err
	}

	return result, nil
}

func (e *ExchangeService) CancelExchange(ctx context.Context, exchangeID uuid.UUID) (ExchangeDTO, error) {
	caller, ok := GetUserIDFromContext(ctx)
	if !ok {
		return ExchangeDTO{}, ErrInvalidCredentials
	}

	var result ExchangeDTO

	err := e.tx.WithTransaction(ctx, func(ctx context.Context) error {
		ex, err := e.exchangeRepo.GetByID(ctx, exchangeID)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return ErrNotFound
			}

			return ErrUnvailible
		}

		if ex.Status() == domain.Completed || ex.Status() == domain.Canceled {
			return ErrInvalidParams
		}

		if ex.FromID() != caller {
			return ErrInvalidCredentials
		}

		updated, err := e.exchangeRepo.UpdateStatus(ctx, exchangeID, domain.Canceled)
		if err != nil {
			return err
		}

		result = mapDomainExchangeToDTO(&updated)

		return nil
	})
	if err != nil {
		return ExchangeDTO{}, err
	}

	return result, nil
}

func (e *ExchangeService) GetUserExchanges(ctx context.Context, userID uuid.UUID, status string) ([]ExchangeDTO, error) {
	exs, err := e.exchangeRepo.GetByUserID(ctx, userID, status)
	if err != nil {
		return nil, ErrUnvailible
	}

	resp := make([]ExchangeDTO, 0, len(exs))
	for _, ex := range exs {
		resp = append(resp, mapDomainExchangeToDTO(ex))
	}

	return resp, nil
}

func mapDomainExchangeToDTO(ex *domain.ExchangeRequest) ExchangeDTO {
	di := ex.DateInfo()

	return ExchangeDTO{
		ID:          ex.ID(),
		OwnedBookID: ex.OwnedBookID(),
		FromID:      ex.FromID(),
		ToID:        ex.ToID(),
		Status:      exchangeStatusToStr(ex.Status()),
		CreatedAt:   di.CreatedAt(),
		ExpiresAt:   di.ExpiresAt(),
		Note:        ex.Note(),
	}
}

func exchangeStatusToStr(s domain.ExchangeStatus) string {
	switch s {
	case domain.Pending:
		return "pending"
	case domain.Accepted:
		return "accepted"
	case domain.Rejected:
		return "rejected"
	case domain.Completed:
		return "completed"
	case domain.Canceled:
		return "canceled"
	default:
		return "pending"
	}
}
