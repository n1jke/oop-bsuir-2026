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
	"golang.org/x/crypto/bcrypt"

	app "github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application"
	mocks "github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application/mock"
)

func TestAuthService_Register(t *testing.T) {
	t.Parallel()

	userID := uuid.New()

	tests := []struct {
		name     string
		username string
		password string
		prepare  func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository)
		wantErr  bool
		check    func(*testing.T, *app.RegisterResponse)
	}{
		{
			name: "success", username: "alice", password: "secret",
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByLogin(gomock.Any(), "alice").Return(app.UserRepoDTO{}, app.ErrNotFound)
				userRepo.EXPECT().Add(gomock.Any(), "alice", gomock.Any()).Return(userID, nil)
			},
			check: func(t *testing.T, resp *app.RegisterResponse) {
				assert.Equal(t, userID, resp.ID)
				assert.NotEmpty(t, resp.Token)
			},
		},
		{
			name: "already exists", username: "alice", password: "secret",
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByLogin(gomock.Any(), "alice").Return(app.UserRepoDTO{ID: userID}, nil)
			},
			wantErr: true,
		},
		{
			name: "get by login error", username: "alice", password: "secret",
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByLogin(gomock.Any(), "alice").Return(app.UserRepoDTO{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "add user error", username: "alice", password: "secret",
			prepare: func(tx *mocks.MockTransactor, userRepo *mocks.MockUserRepository) {
				tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
				userRepo.EXPECT().GetByLogin(gomock.Any(), "alice").Return(app.UserRepoDTO{}, app.ErrNotFound)
				userRepo.EXPECT().Add(gomock.Any(), "alice", gomock.Any()).Return(uuid.Nil, assert.AnError)
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

			tt.prepare(tx, userRepo)

			svc := app.NewAuthService(logger, userRepo, tx, "secret-key", time.Hour)

			resp, err := svc.Register(context.Background(), tt.username, tt.password)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tt.check(t, resp)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	tests := []struct {
		name     string
		username string
		password string
		prepare  func(userRepo *mocks.MockUserRepository)
		wantErr  bool
		check    func(*testing.T, *app.LoginResponse)
	}{
		{
			name: "success", username: "alice", password: "correct-password",
			prepare: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().GetByLogin(gomock.Any(), "alice").Return(app.UserRepoDTO{
					ID: userID, Name: "alice", PasswordHash: string(hash),
				}, nil)
			},
			check: func(t *testing.T, resp *app.LoginResponse) {
				assert.Equal(t, userID, resp.User.ID)
				assert.Equal(t, "alice", resp.User.Name)
				assert.NotEmpty(t, resp.Token)
			},
		},
		{
			name: "user not found", username: "alice", password: "correct-password",
			prepare: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().GetByLogin(gomock.Any(), "alice").Return(app.UserRepoDTO{}, app.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "repo error", username: "alice", password: "correct-password",
			prepare: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().GetByLogin(gomock.Any(), "alice").Return(app.UserRepoDTO{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "wrong password", username: "alice", password: "wrong-password",
			prepare: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().GetByLogin(gomock.Any(), "alice").Return(app.UserRepoDTO{
					ID: userID, Name: "alice", PasswordHash: string(hash),
				}, nil)
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
			userRepo := mocks.NewMockUserRepository(ctrl)
			tx := mocks.NewMockTransactor(ctrl)

			tt.prepare(userRepo)

			svc := app.NewAuthService(logger, userRepo, tx, "secret-key", time.Hour)

			resp, err := svc.Login(context.Background(), tt.username, tt.password)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tt.check(t, resp)
			}
		})
	}
}

func TestUserService_GetAll(t *testing.T) {
	t.Parallel()

	userID := uuid.New()

	tests := []struct {
		name    string
		prepare func(userRepo *mocks.MockUserRepository)
		wantErr bool
		check   func(*testing.T, []*app.UserDTO)
	}{
		{
			name: "success",
			prepare: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().GetAll(gomock.Any()).Return([]*app.UserRepoDTO{
					{ID: userID, Name: "alice", Rating: 4.5},
				}, nil)
			},
			check: func(t *testing.T, users []*app.UserDTO) {
				assert.Len(t, users, 1)
				assert.Equal(t, userID, users[0].ID)
				assert.Equal(t, "alice", users[0].Name)
				assert.Equal(t, 4.5, users[0].Rating)
			},
		},
		{
			name: "repo error",
			prepare: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().GetAll(gomock.Any()).Return(nil, assert.AnError)
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
			userRepo := mocks.NewMockUserRepository(ctrl)

			tt.prepare(userRepo)

			svc := app.NewUserService(logger, userRepo)

			users, err := svc.GetAll(context.Background())
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tt.check(t, users)
			}
		})
	}
}
