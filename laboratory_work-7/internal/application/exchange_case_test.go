package application_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	app "github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application"
	mock "github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application/mock"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/domain"
)

func TestExchangeService_CreateExchange(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	exchangeRepo := mock.NewMockExchangeRepository(ctrl)
	ownedRepo := mock.NewMockOwnedBookRepository(ctrl)
	userRepo := mock.NewMockUserRepository(ctrl)
	tx := mock.NewMockTransactor(ctrl)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	svc := app.NewExchangeService(logger, exchangeRepo, ownedRepo, userRepo, tx)

	fromID := uuid.New()
	toID := uuid.New()
	ownedID := uuid.New()

	t.Run("success", func(t *testing.T) {
		tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })

		userRepo.EXPECT().GetByID(gomock.Any(), toID).Return(app.UserRepoDTO{ID: toID}, nil)

		exchangeRepo.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil)

		ctx := context.WithValue(context.Background(), app.UserIDKeyType{}, fromID)

		dto, err := svc.CreateExchange(ctx, ownedID, toID, time.Now().Add(24*time.Hour), "note")
		require.NoError(t, err)
		assert.Equal(t, fromID, dto.FromID)
		assert.Equal(t, toID, dto.ToID)
		assert.Equal(t, ownedID, dto.OwnedBookID)
	})

	t.Run("to user not found", func(t *testing.T) {
		tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })

		userRepo.EXPECT().GetByID(gomock.Any(), toID).Return(app.UserRepoDTO{}, app.ErrNotFound)

		ctx := context.WithValue(context.Background(), app.UserIDKeyType{}, fromID)

		_, err := svc.CreateExchange(ctx, ownedID, toID, time.Now().Add(24*time.Hour), "note")
		require.Error(t, err)
	})
}

func TestExchangeService_AcceptExchange(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	exchangeRepo := mock.NewMockExchangeRepository(ctrl)
	ownedRepo := mock.NewMockOwnedBookRepository(ctrl)
	userRepo := mock.NewMockUserRepository(ctrl)
	tx := mock.NewMockTransactor(ctrl)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	svc := app.NewExchangeService(logger, exchangeRepo, ownedRepo, userRepo, tx)

	fromID := uuid.New()
	toID := uuid.New()
	ownedID := uuid.New()
	exchID := uuid.New()

	now := time.Now()
	dp := domain.NewDatePolicy(now, now, now.Add(24*time.Hour))
	exReq, _ := domain.CreateExchangeRequest(exchID, ownedID, fromID, toID, domain.Pending, dp, "")

	t.Run("success", func(t *testing.T) {
		tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })

		exchangeRepo.EXPECT().GetByID(gomock.Any(), exchID).Return(*exReq, nil)

		updated, _ := domain.CreateExchangeRequest(exchID, ownedID, fromID, toID, domain.Accepted, dp, "")
		exchangeRepo.EXPECT().UpdateStatus(gomock.Any(), exchID, domain.Accepted).Return(*updated, nil)

		ctx := context.WithValue(context.Background(), app.UserIDKeyType{}, toID)

		dto, err := svc.AcceptExchange(ctx, exchID)
		require.NoError(t, err)
		assert.Equal(t, app.ExchangeDTO{ID: exchID}.ID, dto.ID)
		assert.Equal(t, "accepted", dto.Status)
	})

	t.Run("unauthorized", func(t *testing.T) {
		tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })

		exchangeRepo.EXPECT().GetByID(gomock.Any(), exchID).Return(*exReq, nil)

		ctx := context.WithValue(context.Background(), app.UserIDKeyType{}, uuid.New())

		_, err := svc.AcceptExchange(ctx, exchID)
		require.Error(t, err)
	})
}
