package server

import (
	"context"

	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application"
	codegen "github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/infrastructure/http/codegen"
)

type idKey struct{}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(idKey{}).(uuid.UUID)
	return id, ok
}

type Handler struct {
	authService     *application.AuthService
	userService     *application.UserService
	bookService     *application.BookService
	libraryService  *application.LibraryService
	reviewService   *application.ReviewService
	exchangeService *application.ExchangeService
}

func NewHandler(authService *application.AuthService, userService *application.UserService, bookService *application.BookService,
	libraryService *application.LibraryService, reviewService *application.ReviewService, exchangeService *application.ExchangeService,
) *Handler {
	return &Handler{
		authService:     authService,
		userService:     userService,
		bookService:     bookService,
		libraryService:  libraryService,
		reviewService:   reviewService,
		exchangeService: exchangeService,
	}
}

var _ codegen.ServerInterface = (*Handler)(nil)
