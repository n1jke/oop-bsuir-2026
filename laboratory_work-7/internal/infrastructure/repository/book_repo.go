package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/domain"
)

type BookRepoSQL struct {
	db *pgxpool.Pool
}

const (
	insertBook = `
		INSERT INTO books(id, title, author_id, isbn, description, topics)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, author_id, isbn, description, topics
	`

	selectBookByID = `
		SELECT id, title, author_id, isbn, description, topics FROM books
		WHERE id = $1
	`

	selectBookByISBN = `
		SELECT id, title, author_id, isbn, description, topics FROM books
		WHERE isbn = $1
	`

	selectBookByTitle = `
		SELECT id, title, author_id, isbn, description, topics FROM books
		WHERE title ILIKE '%' || $1 || '%'
	`
)

func (r BookRepoSQL) Add(ctx context.Context, book domain.Book) (domain.Book, error) {
	q := GetQuerier(ctx, r.db)

	var (
		id, authorID uuid.UUID
		title, desc, isbnStr string
		topics []string
	)

	row := q.QueryRow(ctx, insertBook, book.ID(), book.Title(), book.AuthorID(), string(book.ISBN()), book.Description(), book.Topics())
	if err := row.Scan(&id, &title, &authorID, &isbnStr, &desc, &topics); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return domain.Book{}, application.ErrAlreadyExist
		}

		return domain.Book{}, err
	}

	isbn, err := domain.NewISBN(isbnStr)
	if err != nil {
		return domain.Book{}, fmt.Errorf("parse isbn: %w", err)
	}

	created, err := domain.CreateBook(id, title, authorID, isbn, desc, topics...)
	if err != nil {
		return domain.Book{}, err
	}

	return *created, nil
}

func (r BookRepoSQL) GetByID(ctx context.Context, bookID uuid.UUID) (domain.Book, error) {
	q := GetQuerier(ctx, r.db)

	return scanBook(q.QueryRow(ctx, selectBookByID, bookID))
}

func (r BookRepoSQL) GetByISBN(ctx context.Context, isbn string) (domain.Book, error) {
	q := GetQuerier(ctx, r.db)

	return scanBook(q.QueryRow(ctx, selectBookByISBN, isbn))
}

func (r BookRepoSQL) GetByTitle(ctx context.Context, title string) ([]*domain.Book, error) {
	q := GetQuerier(ctx, r.db)

	rows, err := q.Query(ctx, selectBookByTitle, title)
	if err != nil {
		return nil, err
	}

	books := make([]*domain.Book, 0)

	for rows.Next() {
		book, err := scanBookPtr(rows)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func scanBook(row pgx.Row) (domain.Book, error) {
	var (
		id, authorID uuid.UUID
		title, desc, isbnStr string
		topics []string
	)

	if err := row.Scan(&id, &title, &authorID, &isbnStr, &desc, &topics); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Book{}, application.ErrNotFound
		}

		return domain.Book{}, err
	}

	isbn, err := domain.NewISBN(isbnStr)
	if err != nil {
		return domain.Book{}, fmt.Errorf("parse isbn: %w", err)
	}

	book, err := domain.CreateBook(id, title, authorID, isbn, desc, topics...)
	if err != nil {
		return domain.Book{}, err
	}

	return *book, nil
}

func scanBookPtr(row pgx.Row) (*domain.Book, error) {
	var (
		id, authorID uuid.UUID
		title, desc, isbnStr string
		topics []string
	)

	if err := row.Scan(&id, &title, &authorID, &isbnStr, &desc, &topics); err != nil {
		return nil, err
	}

	isbn, err := domain.NewISBN(isbnStr)
	if err != nil {
		return nil, fmt.Errorf("parse isbn: %w", err)
	}

	return domain.CreateBook(id, title, authorID, isbn, desc, topics...)
}

func NewBookRepoSQL(db *pgxpool.Pool) *BookRepoSQL {
	return &BookRepoSQL{db: db}
}
