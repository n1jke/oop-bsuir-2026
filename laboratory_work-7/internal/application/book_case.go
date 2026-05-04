package application

import (
	"context"

	"github.com/google/uuid"
)

type BookService struct{}

func (b *BookService) SearchBook(ctx context.Context, title, topic string) (BookListResponse, error) {
	return BookListResponse{}, nil
}

func (b *BookService) CreateBook(ctx context.Context, title, authorName, isbn, description, topic string) (BookDTO, error) {
	return BookDTO{}, nil
}

func (b *BookService) GetBookByID(ctx context.Context, bookID uuid.UUID) (BookDTO, error) {
	return BookDTO{}, nil
}
