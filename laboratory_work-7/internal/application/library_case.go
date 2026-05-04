package application

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type LibraryService struct {
	logger        *slog.Logger
	ownedBookRepo OwnedBookRepository
	bookRepo      BookRepository
	userRepo      UserRepository
	tx            Transactor
}

func NewLibraryService(logger *slog.Logger, ownedBookRepo OwnedBookRepository, bookRepo BookRepository,
	userRepo UserRepository, tx Transactor,
) *LibraryService {
	logger = logger.With("module", "library-service")

	return &LibraryService{
		logger:        logger,
		ownedBookRepo: ownedBookRepo,
		bookRepo:      bookRepo,
		userRepo:      userRepo,
		tx:            tx,
	}
}

func (l *LibraryService) GetUserBooks(ctx context.Context, userID uuid.UUID) ([]BookDTO, error) {
	var (
		books []BookDTO
		errIn error
	)

	err := l.tx.WithTransaction(ctx, func(context.Context) error {
		books, errIn = l.getUserBooksTx(ctx, userID)
		return errIn
	})
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (l *LibraryService) AddBook(ctx context.Context, userID, bookID uuid.UUID) error {
	return l.tx.WithTransaction(ctx, func(context.Context) error {
		return l.addBookTx(ctx, userID, bookID)
	})
}
