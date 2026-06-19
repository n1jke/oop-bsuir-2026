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

func TestLibraryService_GetUserBooks(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	bookID := uuid.New()
	isbn := validISBN(t)
	book := createTestBook(t, isbn)

	ownedBook, err := domain.CreateOwnedBook(uuid.New(), bookID, userID, domain.Available)
	require.NoError(t, err)

	tests := []struct {
		name    string
		userID  uuid.UUID
		prepare func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, bookRepo *mocks.MockBookRepository,
			ownedRepo *mocks.MockOwnedBookRepository)
		wantErr bool
		check   func(*testing.T, []app.BookDTO)
	}{
		{
			name:   "success",
			userID: userID,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, bookRepo *mocks.MockBookRepository,
				ownedRepo *mocks.MockOwnedBookRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{ID: userID}, nil)
				ownedRepo.EXPECT().GetByUserID(gomock.Any(), userID).Return([]*domain.OwnedBook{ownedBook}, nil)
				bookRepo.EXPECT().GetByID(gomock.Any(), bookID).Return(*book, nil)
			},
			check: func(t *testing.T, books []app.BookDTO) {
				assert.Len(t, books, 1)
				assert.Equal(t, book.Title(), books[0].Title)
			},
		},
		{
			name:   "user not found",
			userID: userID,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, _ *mocks.MockBookRepository,
				_ *mocks.MockOwnedBookRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{}, app.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name:   "user repo error",
			userID: userID,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, _ *mocks.MockBookRepository,
				_ *mocks.MockOwnedBookRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name:   "owned books repo error",
			userID: userID,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, _ *mocks.MockBookRepository,
				ownedRepo *mocks.MockOwnedBookRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{ID: userID}, nil)
				ownedRepo.EXPECT().GetByUserID(gomock.Any(), userID).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			name:   "book not found skip",
			userID: userID,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, bookRepo *mocks.MockBookRepository,
				ownedRepo *mocks.MockOwnedBookRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{ID: userID}, nil)
				ownedRepo.EXPECT().GetByUserID(gomock.Any(), userID).Return([]*domain.OwnedBook{ownedBook}, nil)
				bookRepo.EXPECT().GetByID(gomock.Any(), bookID).Return(domain.Book{}, app.ErrNotFound)
			},
			check: func(t *testing.T, books []app.BookDTO) {
				assert.Len(t, books, 0)
			},
		},
		{
			name:   "book repo error",
			userID: userID,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, bookRepo *mocks.MockBookRepository,
				ownedRepo *mocks.MockOwnedBookRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{ID: userID}, nil)
				ownedRepo.EXPECT().GetByUserID(gomock.Any(), userID).Return([]*domain.OwnedBook{ownedBook}, nil)
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
			userRepo := mocks.NewMockUserRepository(ctrl)
			bookRepo := mocks.NewMockBookRepository(ctrl)
			ownedRepo := mocks.NewMockOwnedBookRepository(ctrl)

			tt.prepare(tx, userRepo, bookRepo, ownedRepo)

			svc := app.NewLibraryService(logger, ownedRepo, bookRepo, userRepo, tx)

			books, err := svc.GetUserBooks(context.Background(), tt.userID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tt.check(t, books)
			}
		})
	}
}

func TestLibraryService_AddBook(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	bookID := uuid.New()

	tests := []struct {
		name    string
		userID  uuid.UUID
		bookID  uuid.UUID
		prepare func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, bookRepo *mocks.MockBookRepository,
			ownedRepo *mocks.MockOwnedBookRepository)
		wantErr bool
	}{
		{
			name:   "success",
			userID: userID, bookID: bookID,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, bookRepo *mocks.MockBookRepository,
				ownedRepo *mocks.MockOwnedBookRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{ID: userID}, nil)
				bookRepo.EXPECT().GetByID(gomock.Any(), bookID).Return(domain.Book{}, nil)
				ownedRepo.EXPECT().GetOwnedBook(gomock.Any(), userID, bookID).Return(domain.OwnedBook{}, app.ErrNotFound)
				ownedRepo.EXPECT().Add(gomock.Any(), gomock.Any()).Return(domain.OwnedBook{}, nil)
			},
		},
		{
			name:   "user not found",
			userID: userID, bookID: bookID,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, _ *mocks.MockBookRepository,
				_ *mocks.MockOwnedBookRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{}, app.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name:   "user repo error",
			userID: userID, bookID: bookID,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, _ *mocks.MockBookRepository,
				_ *mocks.MockOwnedBookRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name:   "book not found",
			userID: userID, bookID: bookID,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, bookRepo *mocks.MockBookRepository,
				_ *mocks.MockOwnedBookRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{ID: userID}, nil)
				bookRepo.EXPECT().GetByID(gomock.Any(), bookID).Return(domain.Book{}, app.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name:   "book repo error",
			userID: userID, bookID: bookID,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, bookRepo *mocks.MockBookRepository,
				_ *mocks.MockOwnedBookRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{ID: userID}, nil)
				bookRepo.EXPECT().GetByID(gomock.Any(), bookID).Return(domain.Book{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name:   "already owned",
			userID: userID, bookID: bookID,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, bookRepo *mocks.MockBookRepository,
				ownedRepo *mocks.MockOwnedBookRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{ID: userID}, nil)
				bookRepo.EXPECT().GetByID(gomock.Any(), bookID).Return(domain.Book{}, nil)
				ownedRepo.EXPECT().GetOwnedBook(gomock.Any(), userID, bookID).Return(domain.OwnedBook{}, nil)
			},
			wantErr: true,
		},
		{
			name:   "get owned book error",
			userID: userID, bookID: bookID,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, bookRepo *mocks.MockBookRepository,
				ownedRepo *mocks.MockOwnedBookRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{ID: userID}, nil)
				bookRepo.EXPECT().GetByID(gomock.Any(), bookID).Return(domain.Book{}, nil)
				ownedRepo.EXPECT().GetOwnedBook(gomock.Any(), userID, bookID).Return(domain.OwnedBook{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name:   "add owned book error",
			userID: userID, bookID: bookID,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, bookRepo *mocks.MockBookRepository,
				ownedRepo *mocks.MockOwnedBookRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{ID: userID}, nil)
				bookRepo.EXPECT().GetByID(gomock.Any(), bookID).Return(domain.Book{}, nil)
				ownedRepo.EXPECT().GetOwnedBook(gomock.Any(), userID, bookID).Return(domain.OwnedBook{}, app.ErrNotFound)
				ownedRepo.EXPECT().Add(gomock.Any(), gomock.Any()).Return(domain.OwnedBook{}, assert.AnError)
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
			userRepo := mocks.NewMockUserRepository(ctrl)
			bookRepo := mocks.NewMockBookRepository(ctrl)
			ownedRepo := mocks.NewMockOwnedBookRepository(ctrl)

			tt.prepare(tx, userRepo, bookRepo, ownedRepo)

			svc := app.NewLibraryService(logger, ownedRepo, bookRepo, userRepo, tx)

			err := svc.AddBook(context.Background(), tt.userID, tt.bookID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
