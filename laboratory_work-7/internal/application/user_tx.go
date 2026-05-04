package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (a *AuthService) registerTx(ctx context.Context, username, password string) (uuid.UUID, error) {
	_, err := a.userRepo.GetByLogin(ctx, username)
	if err == nil {
		return uuid.Nil, ErrUserAlreadyExists
	}

	if !errors.Is(err, ErrNotFound) {
		a.logger.Error("get by id from userRepo", slog.String("username", username), slog.Any("err", err))
		return uuid.Nil, ErrUnvailible
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.logger.Error("bcrypt generate password", slog.Any("err", err))
		return uuid.Nil, fmt.Errorf("bcrypt generate password: %w", err)
	}

	userID, err := a.userRepo.Add(ctx, username, string(hashedPassword))
	if err != nil {
		a.logger.Error("add user to repo", slog.String("username", username), slog.Any("err", err))
		return uuid.Nil, fmt.Errorf("add to repo: %w", err)
	}

	return userID, nil
}
