package application

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/domain"
)

func (l *LibraryService) getUserBooksTx(ctx context.Context, userID uuid.UUID) ([]BookDTO, error) {
	_, err := l.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrUserNotFound
		}

		l.logger.Error("get user by id", slog.String("user_id", userID.String()), slog.Any("err", err))

		return nil, ErrUnvailible
	}

	ownedBooks, err := l.ownedBookRepo.GetByUserID(ctx, userID)
	if err != nil {
		l.logger.Error("get owned books by user", slog.Any("user_id", userID), slog.Any("err", err))
		return nil, ErrUnvailible
	}

	books := make([]BookDTO, 0, len(ownedBooks))
	for _, ob := range ownedBooks {
		book, err := l.bookRepo.GetByID(ctx, ob.BookID())
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				continue
			}

			l.logger.Error("get book by id", slog.Any("book_id", ob.BookID()), slog.Any("err", err))

			return nil, ErrUnvailible
		}

		books = append(books, BookDTO{
			ID:          book.ID(),
			Title:       book.Title(),
			AuthorName:  "",
			ISBN:        string(book.ISBN()),
			Description: book.Description(),
			Topics:      book.Topics(),
		})
	}

	return books, nil
}

func (l *LibraryService) addBookTx(ctx context.Context, userID, bookID uuid.UUID) error {
	_, err := l.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrUserNotFound
		}

		l.logger.Error("get user by id", slog.Any("user_id", userID), slog.Any("err", err))

		return ErrUnvailible
	}

	_, err = l.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrBookNotFound
		}

		l.logger.Error("get book by id", slog.Any("book_id", bookID), slog.Any("err", err))

		return ErrUnvailible
	}

	_, err = l.ownedBookRepo.GetOwnedBook(ctx, userID, bookID)
	if err == nil {
		return ErrAlreadyExist
	}

	if !errors.Is(err, ErrNotFound) {
		l.logger.Error("get owned book", slog.Any("err", err))
		return ErrUnvailible
	}

	ownedBook, err := domain.NewOwnedBook(bookID, userID)
	if err != nil {
		return ErrInvalidParams
	}

	_, err = l.ownedBookRepo.Add(ctx, *ownedBook)
	if err != nil {
		l.logger.Error("add owned book to repo", slog.Any("err", err))
		return ErrUnvailible
	}

	return nil
}
