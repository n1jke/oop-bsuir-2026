package application

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type BookService struct {
	logger   *slog.Logger
	bookRepo BookRepository
	tx       Transactor
}

func NewBookService(logger *slog.Logger, bookRepo BookRepository, tx Transactor) *BookService {
	logger = logger.With("module", "book-service")

	return &BookService{
		logger:   logger,
		bookRepo: bookRepo,
		tx:       tx,
	}
}

func (b *BookService) SearchBook(ctx context.Context, title, topic string) (BookListResponse, error) {
	var (
		result BookListResponse
		errIn  error
	)

	err := b.tx.WithTransaction(ctx, func(context.Context) error {
		result, errIn = b.searchBookTx(ctx, title)
		return errIn
	})
	if err != nil {
		return BookListResponse{}, err
	}

	return result, nil
}

func (b *BookService) CreateBook(ctx context.Context, title, authorName, isbn, description, topic string) (BookDTO, error) {
	var (
		book  BookDTO
		errIn error
	)

	err := b.tx.WithTransaction(ctx, func(context.Context) error {
		book, errIn = b.createBookTx(ctx, title, authorName, isbn, description, topic)
		return errIn
	})
	if err != nil {
		return BookDTO{}, err
	}

	return book, nil
}

func (b *BookService) GetBookByID(ctx context.Context, bookID uuid.UUID) (BookDTO, error) {
	var (
		book  BookDTO
		errIn error
	)

	err := b.tx.WithTransaction(ctx, func(context.Context) error {
		book, errIn = b.getBookByIDTx(ctx, bookID)
		return errIn
	})
	if err != nil {
		return BookDTO{}, err
	}

	return book, nil
}
