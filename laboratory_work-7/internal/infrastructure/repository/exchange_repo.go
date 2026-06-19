package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/domain"
)

type ExchangeRepoSQL struct {
	db *pgxpool.Pool
}

const (
	insertExchange = `
		INSERT INTO exchanges(id, owned_book_id, from_id, to_id, status, note, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	selectExchangeByID = `
		SELECT id, owned_book_id, from_id, to_id, status, note, created_at, updated_at, expires_at
		FROM exchanges WHERE id = $1
	`

	selectExchangesByUser = `
		SELECT id, owned_book_id, from_id, to_id, status, note, created_at, updated_at, expires_at
		FROM exchanges WHERE from_id = $1 OR to_id = $1
		ORDER BY created_at DESC
	`

	selectExchangesByUserWithStatus = `
		SELECT id, owned_book_id, from_id, to_id, status, note, created_at, updated_at, expires_at
		FROM exchanges WHERE (from_id = $1 OR to_id = $1) AND status = $2
		ORDER BY created_at DESC
	`

	selectExchangesByOwnedBook = `
		SELECT id, owned_book_id, from_id, to_id, status, note, created_at, updated_at, expires_at
		FROM exchanges WHERE owned_book_id = $1
	`

	updateExchangeStatus = `
		UPDATE exchanges SET status = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, owned_book_id, from_id, to_id, status, note, created_at, updated_at, expires_at
	`
)

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
		return ""
	}
}

func strToExchangeStatus(s string) (domain.ExchangeStatus, error) {
	switch s {
	case "pending":
		return domain.Pending, nil
	case "accepted":
		return domain.Accepted, nil
	case "rejected":
		return domain.Rejected, nil
	case "completed":
		return domain.Completed, nil
	case "canceled":
		return domain.Canceled, nil
	default:
		return 0, application.ErrInvalidParams
	}
}

func (r ExchangeRepoSQL) Add(ctx context.Context, exchange *domain.ExchangeRequest) error {
	q := GetQuerier(ctx, r.db)

	info := exchange.DateInfo()
	status := exchangeStatusToStr(exchange.Status())

	_, err := q.Exec(ctx, insertExchange,
		exchange.ID(), exchange.OwnedBookID(), exchange.FromID(), exchange.ToID(),
		status, exchange.Note(), info.CreatedAt(), info.ExpiresAt(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r ExchangeRepoSQL) GetByID(ctx context.Context, exchangeID uuid.UUID) (domain.ExchangeRequest, error) {
	q := GetQuerier(ctx, r.db)

	return scanExchange(q.QueryRow(ctx, selectExchangeByID, exchangeID))
}

func (r ExchangeRepoSQL) GetByUserID(ctx context.Context, userID uuid.UUID, status string) ([]*domain.ExchangeRequest, error) {
	q := GetQuerier(ctx, r.db)

	var (
		rows pgx.Rows
		err  error
	)

	if status == "" {
		rows, err = q.Query(ctx, selectExchangesByUser, userID)
	} else {
		rows, err = q.Query(ctx, selectExchangesByUserWithStatus, userID, status)
	}

	if err != nil {
		return nil, err
	}

	exchanges := make([]*domain.ExchangeRequest, 0)

	for rows.Next() {
		ex, err := scanExchangePtr(rows)
		if err != nil {
			return nil, err
		}

		exchanges = append(exchanges, ex)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return exchanges, nil
}

func (r ExchangeRepoSQL) GetByOwnedBook(ctx context.Context, ownedBookID uuid.UUID) ([]*domain.ExchangeRequest, error) {
	q := GetQuerier(ctx, r.db)

	rows, err := q.Query(ctx, selectExchangesByOwnedBook, ownedBookID)
	if err != nil {
		return nil, err
	}

	exchanges := make([]*domain.ExchangeRequest, 0)

	for rows.Next() {
		ex, err := scanExchangePtr(rows)
		if err != nil {
			return nil, err
		}

		exchanges = append(exchanges, ex)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return exchanges, nil
}

func (r ExchangeRepoSQL) UpdateStatus(ctx context.Context, exchangeID uuid.UUID, status domain.ExchangeStatus,
) (domain.ExchangeRequest, error) {
	q := GetQuerier(ctx, r.db)

	statusStr := exchangeStatusToStr(status)

	return scanExchange(q.QueryRow(ctx, updateExchangeStatus, exchangeID, statusStr))
}

func scanExchange(row pgx.Row) (domain.ExchangeRequest, error) {
	var (
		id, ownedBookID, fromID, toID   uuid.UUID
		statusStr, note                 string
		createdAt, updatedAt, expiresAt time.Time
	)

	if err := row.Scan(&id, &ownedBookID, &fromID, &toID, &statusStr, &note, &createdAt, &updatedAt, &expiresAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ExchangeRequest{}, application.ErrNotFound
		}

		return domain.ExchangeRequest{}, err
	}

	status, err := strToExchangeStatus(statusStr)
	if err != nil {
		return domain.ExchangeRequest{}, err
	}

	dateInfo := domain.NewDatePolicy(createdAt, updatedAt, expiresAt)

	ex, err := domain.CreateExchangeRequest(id, ownedBookID, fromID, toID, status, dateInfo, note)
	if err != nil {
		return domain.ExchangeRequest{}, err
	}

	return *ex, nil
}

func scanExchangePtr(row pgx.Row) (*domain.ExchangeRequest, error) {
	var (
		id, ownedBookID, fromID, toID   uuid.UUID
		statusStr, note                 string
		createdAt, updatedAt, expiresAt time.Time
	)

	if err := row.Scan(&id, &ownedBookID, &fromID, &toID, &statusStr, &note, &createdAt, &updatedAt, &expiresAt); err != nil {
		return nil, err
	}

	status, err := strToExchangeStatus(statusStr)
	if err != nil {
		return nil, err
	}

	dateInfo := domain.NewDatePolicy(createdAt, updatedAt, expiresAt)

	return domain.CreateExchangeRequest(id, ownedBookID, fromID, toID, status, dateInfo, note)
}

func NewExchangeRepoSQL(db *pgxpool.Pool) *ExchangeRepoSQL {
	return &ExchangeRepoSQL{db: db}
}
