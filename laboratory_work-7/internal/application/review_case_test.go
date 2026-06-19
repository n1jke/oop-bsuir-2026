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

func TestReviewService_AddBookReview(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	bookID := uuid.New()

	tests := []struct {
		name    string
		userID  uuid.UUID
		bookID  uuid.UUID
		mark    uint
		report  string
		prepare func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, bookRepo *mocks.MockBookRepository,
			reviewRepo *mocks.MockReviewRepository)
		wantErr bool
		check   func(*testing.T, *app.BookReviewDTO)
	}{
		{
			name: "success", userID: userID, bookID: bookID, mark: 5, report: "good",
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, bookRepo *mocks.MockBookRepository,
				reviewRepo *mocks.MockReviewRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{ID: userID}, nil)
				bookRepo.EXPECT().GetByID(gomock.Any(), bookID).Return(createBook(), nil)
				reviewRepo.EXPECT().AddBookReview(gomock.Any(), gomock.Any()).Return(nil)
			},
			check: func(t *testing.T, dto *app.BookReviewDTO) {
				assert.Equal(t, userID, dto.FromID)
				assert.Equal(t, bookID, dto.BookID)
				assert.Equal(t, uint(5), dto.Mark)
			},
		},
		{
			name: "user not found", userID: userID, bookID: bookID, mark: 5,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, _ *mocks.MockBookRepository, _ *mocks.MockReviewRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{}, app.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "user repo error", userID: userID, bookID: bookID, mark: 5,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, _ *mocks.MockBookRepository, _ *mocks.MockReviewRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "book not found", userID: userID, bookID: bookID, mark: 5,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, bookRepo *mocks.MockBookRepository,
				_ *mocks.MockReviewRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{ID: userID}, nil)
				bookRepo.EXPECT().GetByID(gomock.Any(), bookID).Return(domain.Book{}, app.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "invalid mark", userID: userID, bookID: bookID, mark: 11,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, bookRepo *mocks.MockBookRepository,
				_ *mocks.MockReviewRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{ID: userID}, nil)
				bookRepo.EXPECT().GetByID(gomock.Any(), bookID).Return(createBook(), nil)
			},
			wantErr: true,
		},
		{
			name: "add review error", userID: userID, bookID: bookID, mark: 5,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, bookRepo *mocks.MockBookRepository,
				reviewRepo *mocks.MockReviewRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(app.UserRepoDTO{ID: userID}, nil)
				bookRepo.EXPECT().GetByID(gomock.Any(), bookID).Return(createBook(), nil)
				reviewRepo.EXPECT().AddBookReview(gomock.Any(), gomock.Any()).Return(assert.AnError)
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
			reviewRepo := mocks.NewMockReviewRepository(ctrl)

			tt.prepare(tx, userRepo, bookRepo, reviewRepo)

			svc := app.NewReviewService(logger, reviewRepo, bookRepo, userRepo, tx)

			dto, err := svc.AddBookReview(context.Background(), tt.userID, tt.bookID, tt.mark, tt.report)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tt.check(t, dto)
			}
		})
	}
}

func TestReviewService_AddUserReview(t *testing.T) {
	t.Parallel()

	fromID := uuid.New()
	toID := uuid.New()

	tests := []struct {
		name    string
		fromID  uuid.UUID
		toID    uuid.UUID
		mark    uint
		report  string
		prepare func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, bookRepo *mocks.MockBookRepository,
			reviewRepo *mocks.MockReviewRepository)
		wantErr bool
		check   func(*testing.T, *app.UserReviewDTO)
	}{
		{
			name: "success", fromID: fromID, toID: toID, mark: 7, report: "nice",
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, _ *mocks.MockBookRepository,
				reviewRepo *mocks.MockReviewRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), fromID).Return(app.UserRepoDTO{ID: fromID}, nil)
				userRepo.EXPECT().GetByID(gomock.Any(), toID).Return(app.UserRepoDTO{ID: toID}, nil)
				reviewRepo.EXPECT().AddUserReview(gomock.Any(), gomock.Any()).Return(nil)
			},
			check: func(t *testing.T, dto *app.UserReviewDTO) {
				assert.Equal(t, fromID, dto.FromID)
				assert.Equal(t, toID, dto.ToID)
				assert.Equal(t, uint(7), dto.Mark)
			},
		},
		{
			name: "from user not found", fromID: fromID, toID: toID, mark: 7,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, _ *mocks.MockBookRepository, _ *mocks.MockReviewRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), fromID).Return(app.UserRepoDTO{}, app.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "from user repo error", fromID: fromID, toID: toID, mark: 7,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, _ *mocks.MockBookRepository, _ *mocks.MockReviewRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), fromID).Return(app.UserRepoDTO{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "to user not found", fromID: fromID, toID: toID, mark: 7,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, _ *mocks.MockBookRepository, _ *mocks.MockReviewRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), fromID).Return(app.UserRepoDTO{ID: fromID}, nil)
				userRepo.EXPECT().GetByID(gomock.Any(), toID).Return(app.UserRepoDTO{}, app.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "invalid mark", fromID: fromID, toID: toID, mark: 11,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, _ *mocks.MockBookRepository, _ *mocks.MockReviewRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), fromID).Return(app.UserRepoDTO{ID: fromID}, nil)
				userRepo.EXPECT().GetByID(gomock.Any(), toID).Return(app.UserRepoDTO{ID: toID}, nil)
			},
			wantErr: true,
		},
		{
			name: "self review", fromID: fromID, toID: fromID, mark: 5,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, _ *mocks.MockBookRepository, _ *mocks.MockReviewRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), fromID).Return(app.UserRepoDTO{ID: fromID}, nil)
				userRepo.EXPECT().GetByID(gomock.Any(), fromID).Return(app.UserRepoDTO{ID: fromID}, nil)
			},
			wantErr: true,
		},
		{
			name: "add review error", fromID: fromID, toID: toID, mark: 7,
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository, _ *mocks.MockBookRepository,
				reviewRepo *mocks.MockReviewRepository,
			) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByID(gomock.Any(), fromID).Return(app.UserRepoDTO{ID: fromID}, nil)
				userRepo.EXPECT().GetByID(gomock.Any(), toID).Return(app.UserRepoDTO{ID: toID}, nil)
				reviewRepo.EXPECT().AddUserReview(gomock.Any(), gomock.Any()).Return(assert.AnError)
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
			reviewRepo := mocks.NewMockReviewRepository(ctrl)

			tt.prepare(tx, userRepo, bookRepo, reviewRepo)

			svc := app.NewReviewService(logger, reviewRepo, bookRepo, userRepo, tx)

			dto, err := svc.AddUserReview(context.Background(), tt.fromID, tt.toID, tt.mark, tt.report)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tt.check(t, dto)
			}
		})
	}
}

func createBook() domain.Book {
	isbn, _ := domain.NewISBN("978-0-123456-47-2")
	book, _ := domain.CreateBook(uuid.New(), "Test Book", uuid.Nil, isbn, "Description", "topic1")

	return *book
}
