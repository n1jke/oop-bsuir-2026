package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/domain"
)

type OwnedBookRepoSQL struct {
	db *pgxpool.Pool
}

const (
	insertOwnedBook = `
		INSERT INTO owned_books(id, book_id, owner_id, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, book_id, owner_id, status
	`

	selectOwnedBooksByUser = `
		SELECT id, book_id, owner_id, status FROM owned_books
		WHERE owner_id = $1
	`

	selectOwnedBook = `
		SELECT id, book_id, owner_id, status FROM owned_books
		WHERE owner_id = $1 AND book_id = $2
	`

	updateOwnedBookStatus = `
		UPDATE owned_books SET status = $2
		WHERE id = $1
		RETURNING id, book_id, owner_id, status
	`

	uniqueViolation = "23505"
)

func (r OwnedBookRepoSQL) Add(ctx context.Context, ownedBook domain.OwnedBook) (domain.OwnedBook, error) {
	q := GetQuerier(ctx, r.db)

	var (
		id, bookID, ownerID uuid.UUID
		status              domain.OwnedBookStatus
	)

	row := q.QueryRow(ctx, insertOwnedBook, ownedBook.ID(), ownedBook.BookID(), ownedBook.OwnerID(), ownedBook.Status())
	if err := row.Scan(&id, &bookID, &ownerID, &status); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return domain.OwnedBook{}, application.ErrAlreadyExist
		}

		return domain.OwnedBook{}, err
	}

	created, err := domain.CreateOwnedBook(id, bookID, ownerID, status)
	if err != nil {
		return domain.OwnedBook{}, err
	}

	return *created, nil
}

func (r OwnedBookRepoSQL) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.OwnedBook, error) {
	q := GetQuerier(ctx, r.db)

	rows, err := q.Query(ctx, selectOwnedBooksByUser, userID)
	if err != nil {
		return nil, err
	}

	books := make([]*domain.OwnedBook, 0)

	for rows.Next() {
		var (
			id, bookID, ownerID uuid.UUID
			status              domain.OwnedBookStatus
		)

		if err := rows.Scan(&id, &bookID, &ownerID, &status); err != nil {
			return nil, err
		}

		ob, err := domain.CreateOwnedBook(id, bookID, ownerID, status)
		if err != nil {
			return nil, err
		}

		books = append(books, ob)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r OwnedBookRepoSQL) GetOwnedBook(ctx context.Context, userID, bookID uuid.UUID) (domain.OwnedBook, error) {
	q := GetQuerier(ctx, r.db)

	var (
		id, obID, ownerID uuid.UUID
		status            domain.OwnedBookStatus
	)

	row := q.QueryRow(ctx, selectOwnedBook, userID, bookID)
	if err := row.Scan(&id, &obID, &ownerID, &status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.OwnedBook{}, application.ErrNotFound
		}

		return domain.OwnedBook{}, err
	}

	ob, err := domain.CreateOwnedBook(id, obID, ownerID, status)
	if err != nil {
		return domain.OwnedBook{}, err
	}

	return *ob, nil
}

func (r OwnedBookRepoSQL) UpdateStatus(ctx context.Context, ownedBookID uuid.UUID, status domain.OwnedBookStatus,
) (domain.OwnedBook, error) {
	q := GetQuerier(ctx, r.db)

	var (
		id, bookID, ownerID uuid.UUID
		newStatus           domain.OwnedBookStatus
	)

	row := q.QueryRow(ctx, updateOwnedBookStatus, ownedBookID, status)
	if err := row.Scan(&id, &bookID, &ownerID, &newStatus); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.OwnedBook{}, application.ErrNotFound
		}

		return domain.OwnedBook{}, err
	}

	ob, err := domain.CreateOwnedBook(id, bookID, ownerID, newStatus)
	if err != nil {
		return domain.OwnedBook{}, err
	}

	return *ob, nil
}

func NewOwnedBookRepoSQL(db *pgxpool.Pool) *OwnedBookRepoSQL {
	return &OwnedBookRepoSQL{db: db}
}
