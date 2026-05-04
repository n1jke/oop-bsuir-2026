package application

import (
	"context"

	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/domain"
)

//go:generate mockgen -source=interfaces.go -destination=mock/mock.go -package=mocks

type Transactor interface {
	WithTransaction(context.Context, func(context.Context) error) error
}

type BookRepository interface {
	Add(ctx context.Context, book domain.Book) (domain.Book, error)
	GetByTitle(ctx context.Context, title string) ([]domain.Book, error)
	GetByID(ctx context.Context, bookID uuid.UUID) (domain.Book, error)
}

type ExchangeRepository interface {
	Add(ctx context.Context, exchange domain.ExchangeRequest) error
	GetByID(ctx context.Context, exchangeID uuid.UUID) (domain.ExchangeRequest, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, status string) ([]domain.ExchangeRequest, error)
	GetByOwnedBook(ctx context.Context, ownedBookID uuid.UUID) ([]domain.ExchangeRequest, error)
	UpdateStatus(ctx context.Context, exchangeID uuid.UUID, status domain.ExchangeStatus) (domain.ExchangeRequest, error)
}

type OwnedBookRepository interface {
	Add(ctx context.Context, ownedBook domain.OwnedBook) (domain.OwnedBook, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.OwnedBook, error)
	GetOwnedBook(ctx context.Context, userID, bookID uuid.UUID) (domain.OwnedBook, error)
	UpdateStatus(ctx context.Context, ownedBookID uuid.UUID, status domain.OwnedBookStatus) (domain.OwnedBook, error)
}

type ReviewRepository interface {
	AddBookReview(ctx context.Context, review domain.BookReview) error
	AddUserReview(ctx context.Context, review domain.UserReview) error
	GetBookReviews(ctx context.Context, bookID uuid.UUID) ([]domain.BookReview, error)
	GetUserReviews(ctx context.Context, userID uuid.UUID) ([]domain.UserReview, error)
}

type UserRepository interface {
	Add(ctx context.Context, login, password string) (uuid.UUID, error)
	GetByLogin(ctx context.Context, login string) (UserRepoDTO, error)
	GetByID(ctx context.Context, userID uuid.UUID) (UserRepoDTO, error)
	GetAll(ctx context.Context) ([]UserRepoDTO, error)
}
