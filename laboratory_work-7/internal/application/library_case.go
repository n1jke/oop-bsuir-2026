package application

import (
	"context"

	"github.com/google/uuid"
)

type LibraryService struct{}

func (l *LibraryService) GetUserBooks(ctx context.Context, userID uuid.UUID) ([]BookDTO, error) {
	return nil, nil
}

func (l *LibraryService) AddBook(ctx context.Context, userID, bookID uuid.UUID) error {
	return nil
}
