package application_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	app "github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application"
	mocks "github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application/mock"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/domain"
)

func TestBookService_SearchBook(t *testing.T) {
	t.Parallel()

	isbn := validISBN(t)
	book := createTestBook(t, isbn)

	tests := []struct {
		name    string
		title   string
		prepare func(tx *mocks.MockTransactor, bookRepo *mocks.MockBookRepository)
		wantErr bool
		check   func(*testing.T, app.BookListResponse)
	}{
		{
			name:  "success",
			title: "Test",
			prepare: func(tx *mocks.MockTransactor, bookRepo *mocks.MockBookRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				bookRepo.EXPECT().GetByTitle(gomock.Any(), "Test").Return([]*domain.Book{book}, nil)
			},
			check: func(t *testing.T, resp app.BookListResponse) {
				assert.Len(t, resp.Books, 1)
				assert.Equal(t, "Test Book", resp.Books[0].Title)
				assert.Equal(t, "978-0-123456-47-2", resp.Books[0].ISBN)
				assert.Equal(t, 1, resp.Total)
			},
		},
		{
			name:  "empty result",
			title: "Unknown",
			prepare: func(tx *mocks.MockTransactor, bookRepo *mocks.MockBookRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				bookRepo.EXPECT().GetByTitle(gomock.Any(), "Unknown").Return([]*domain.Book{}, nil)
			},
			check: func(t *testing.T, resp app.BookListResponse) {
				assert.Len(t, resp.Books, 0)
				assert.Equal(t, 0, resp.Total)
			},
		},
		{
			name:  "repo error",
			title: "Test",
			prepare: func(tx *mocks.MockTransactor, bookRepo *mocks.MockBookRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				bookRepo.EXPECT().GetByTitle(gomock.Any(), "Test").Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
			tx := mocks.NewMockTransactor(ctrl)
			bookRepo := mocks.NewMockBookRepository(ctrl)

			tt.prepare(tx, bookRepo)

			svc := app.NewBookService(logger, bookRepo, tx)

			resp, err := svc.SearchBook(context.Background(), tt.title, "")
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tt.check(t, resp)
			}
		})
	}
}

func TestBookService_CreateBook(t *testing.T) {
	t.Parallel()

	isbnStr := string(validISBN(t))
	savedBook := func() domain.Book {
		isbn, err := domain.NewISBN(isbnStr)
		require.NoError(t, err)
		b, err := domain.CreateBook(uuid.New(), "New Book", uuid.Nil, isbn, "Desc", "fiction")
		require.NoError(t, err)

		return *b
	}()

	tests := []struct {
		name        string
		title       string
		authorName  string
		isbn        string
		description string
		topic       string
		prepare     func(tx *mocks.MockTransactor, bookRepo *mocks.MockBookRepository)
		wantErr     bool
		check       func(*testing.T, app.BookDTO)
	}{
		{
			name:  "success",
			title: "New Book", isbn: isbnStr, description: "Desc", topic: "fiction",
			prepare: func(tx *mocks.MockTransactor, bookRepo *mocks.MockBookRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				bookRepo.EXPECT().GetByISBN(gomock.Any(), isbnStr).Return(domain.Book{}, app.ErrNotFound)
				bookRepo.EXPECT().Add(gomock.Any(), gomock.Any()).Return(savedBook, nil)
			},
			check: func(t *testing.T, dto app.BookDTO) {
				assert.Equal(t, savedBook.Title(), dto.Title)
				assert.Equal(t, string(savedBook.ISBN()), dto.ISBN)
			},
		},
		{
			name:  "invalid isbn",
			title: "New Book", isbn: "invalid", description: "Desc", topic: "fiction",
			prepare: func(tx *mocks.MockTransactor, _ *mocks.MockBookRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
			},
			wantErr: true,
		},
		{
			name:  "duplicate isbn",
			title: "New Book", isbn: isbnStr, description: "Desc", topic: "fiction",
			prepare: func(tx *mocks.MockTransactor, bookRepo *mocks.MockBookRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				bookRepo.EXPECT().GetByISBN(gomock.Any(), isbnStr).Return(domain.Book{}, nil)
			},
			wantErr: true,
		},
		{
			name:  "get by isbn unvailible",
			title: "New Book", isbn: isbnStr, description: "Desc", topic: "fiction",
			prepare: func(tx *mocks.MockTransactor, bookRepo *mocks.MockBookRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				bookRepo.EXPECT().GetByISBN(gomock.Any(), isbnStr).Return(domain.Book{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name:  "add repo error",
			title: "New Book", isbn: isbnStr, description: "Desc", topic: "fiction",
			prepare: func(tx *mocks.MockTransactor, bookRepo *mocks.MockBookRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				bookRepo.EXPECT().GetByISBN(gomock.Any(), isbnStr).Return(domain.Book{}, app.ErrNotFound)
				bookRepo.EXPECT().Add(gomock.Any(), gomock.Any()).Return(domain.Book{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name:  "description too long",
			title: "New Book", isbn: isbnStr, description: longDesc(256), topic: "fiction",
			prepare: func(tx *mocks.MockTransactor, bookRepo *mocks.MockBookRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				bookRepo.EXPECT().GetByISBN(gomock.Any(), isbnStr).Return(domain.Book{}, app.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
			tx := mocks.NewMockTransactor(ctrl)
			bookRepo := mocks.NewMockBookRepository(ctrl)

			tt.prepare(tx, bookRepo)

			svc := app.NewBookService(logger, bookRepo, tx)

			dto, err := svc.CreateBook(context.Background(), tt.title, tt.authorName, tt.isbn, tt.description, tt.topic)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tt.check(t, dto)
			}
		})
	}
}

func TestBookService_GetBookByID(t *testing.T) {
	t.Parallel()

	isbn := validISBN(t)
	book := createTestBook(t, isbn)
	bookID := uuid.New()

	tests := []struct {
		name    string
		bookID  uuid.UUID
		prepare func(tx *mocks.MockTransactor, bookRepo *mocks.MockBookRepository)
		wantErr bool
		check   func(*testing.T, app.BookDTO)
	}{
		{
			name:   "success",
			bookID: bookID,
			prepare: func(tx *mocks.MockTransactor, bookRepo *mocks.MockBookRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				bookRepo.EXPECT().GetByID(gomock.Any(), bookID).Return(*book, nil)
			},
			check: func(t *testing.T, dto app.BookDTO) {
				assert.Equal(t, book.ID(), dto.ID)
				assert.Equal(t, book.Title(), dto.Title)
			},
		},
		{
			name:   "not found",
			bookID: bookID,
			prepare: func(tx *mocks.MockTransactor, bookRepo *mocks.MockBookRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				bookRepo.EXPECT().GetByID(gomock.Any(), bookID).Return(domain.Book{}, app.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name:   "repo error",
			bookID: bookID,
			prepare: func(tx *mocks.MockTransactor, bookRepo *mocks.MockBookRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				bookRepo.EXPECT().GetByID(gomock.Any(), bookID).Return(domain.Book{}, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
			tx := mocks.NewMockTransactor(ctrl)
			bookRepo := mocks.NewMockBookRepository(ctrl)

			tt.prepare(tx, bookRepo)

			svc := app.NewBookService(logger, bookRepo, tx)

			dto, err := svc.GetBookByID(context.Background(), tt.bookID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tt.check(t, dto)
			}
		})
	}
}

func validISBN(t *testing.T) domain.ISBN {
	t.Helper()

	isbn, err := domain.NewISBN("978-0-123456-47-2")
	require.NoError(t, err)

	return isbn
}

func createTestBook(t *testing.T, isbn domain.ISBN) *domain.Book {
	t.Helper()

	book, err := domain.CreateBook(uuid.New(), "Test Book", uuid.Nil, isbn, "Description", "topic1")
	require.NoError(t, err)

	return book
}

func longDesc(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = 'a'
	}

	return string(b)
}
