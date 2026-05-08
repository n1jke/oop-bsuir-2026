package application

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/domain"
)

func (b *BookService) searchBookTx(ctx context.Context, title string) (BookListResponse, error) {
	books, err := b.bookRepo.GetByTitle(ctx, title)
	if err != nil {
		b.logger.Error("get books by title", slog.String("title", title), slog.Any("err", err))
		return BookListResponse{}, ErrUnvailible
	}

	dtos := make([]BookDTO, 0, len(books))
	for _, book := range books {
		dtos = append(dtos, BookDTO{
			ID:          book.ID(),
			Title:       book.Title(),
			AuthorName:  "",
			ISBN:        string(book.ISBN()),
			Description: book.Description(),
			Topics:      book.Topics(),
		})
	}

	return BookListResponse{
		Books: dtos,
		Total: len(dtos),
	}, nil
}

func (b *BookService) createBookTx(ctx context.Context, title, authorName, isbn, description, topic string) (BookDTO, error) {
	isbnObj, err := domain.NewISBN(isbn)
	if err != nil {
		return BookDTO{}, ErrInvalidParams
	}

	_, err = b.bookRepo.GetByISBN(ctx, isbn)
	if err == nil {
		return BookDTO{}, ErrAlreadyExist
	}

	if !errors.Is(err, ErrNotFound) {
		b.logger.Error("get book by isbn", slog.String("isbn", isbn), slog.Any("err", err))
		return BookDTO{}, ErrUnvailible
	}

	book, err := domain.NewBook(title, uuid.Nil, isbnObj, description, topic)
	if err != nil {
		return BookDTO{}, ErrInvalidParams
	}

	saved, err := b.bookRepo.Add(ctx, book)
	if err != nil {
		b.logger.Error("add book to repo", slog.Any("err", err))
		return BookDTO{}, ErrUnvailible
	}

	return BookDTO{
		ID:          saved.ID(),
		Title:       saved.Title(),
		AuthorName:  authorName,
		ISBN:        string(saved.ISBN()),
		Description: saved.Description(),
		Topics:      saved.Topics(),
	}, nil
}

func (b *BookService) getBookByIDTx(ctx context.Context, bookID uuid.UUID) (BookDTO, error) {
	book, err := b.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return BookDTO{}, ErrBookNotFound
		}

		b.logger.Error("get book by id", slog.String("book_id", bookID.String()), slog.Any("err", err))

		return BookDTO{}, ErrUnvailible
	}

	return BookDTO{
		ID:          book.ID(),
		Title:       book.Title(),
		AuthorName:  "",
		ISBN:        string(book.ISBN()),
		Description: book.Description(),
		Topics:      book.Topics(),
	}, nil
}
