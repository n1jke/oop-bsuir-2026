package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	logger         *slog.Logger
	userRepo       UserRepository
	tx             Transactor
	secretKey      []byte
	expirationTime time.Duration
}

func NewAuthService(logger *slog.Logger, repo UserRepository, tx Transactor, secret string, expTime time.Duration) *AuthService {
	logger = logger.With("module", "auth-service")

	return &AuthService{
		logger:         logger,
		userRepo:       repo,
		tx:             tx,
		secretKey:      []byte(secret),
		expirationTime: expTime,
	}
}

func (a *AuthService) Register(ctx context.Context, username, password string) (*RegisterResponse, error) {
	var (
		errIn  error
		userID uuid.UUID
	)

	err := a.tx.WithTransaction(ctx, func(ctx context.Context) error {
		userID, errIn = a.registerTx(ctx, username, password)
		return errIn
	})
	if err != nil {
		return nil, err
	}

	token, err := a.generateJWT(userID)
	if err != nil {
		a.logger.Error("generating JWT", slog.Any("err", err))
		return nil, fmt.Errorf("generating JWT: %w", err)
	}

	return &RegisterResponse{
		ID:    userID,
		Token: token,
	}, nil
}

func (a *AuthService) Login(ctx context.Context, username, password string) (*LoginResponse, error) {
	userDTO, err := a.userRepo.GetByLogin(ctx, username)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrInvalidCredentials
		}

		a.logger.Error("get by id from userRepo", slog.String("username", username), slog.Any("err", err))

		return nil, ErrUnvailible
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDTO.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := a.generateJWT(userDTO.ID)
	if err != nil {
		a.logger.Error("generating JWT", slog.Any("err", err))
		return nil, fmt.Errorf("generating JWT: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User: UserDTO{
			ID:     userDTO.ID,
			Name:   userDTO.Name,
			Rating: userDTO.Rating,
		},
	}, nil
}

type UserService struct {
	logger   *slog.Logger
	userRepo UserRepository
}

func NewUserService(logger *slog.Logger, userRepo UserRepository) *UserService {
	logger = logger.With("module", "user-service")

	return &UserService{
		logger:   logger,
		userRepo: userRepo,
	}
}

func (u *UserService) GetAll(ctx context.Context) ([]*UserDTO, error) {
	users, err := u.userRepo.GetAll(ctx)
	if err != nil {
		u.logger.Error("get users from userRepo", slog.Any("err", err))
		return nil, ErrUnvailible
	}

	dtoList := make([]*UserDTO, 0, len(users))
	for _, user := range users {
		dtoList = append(dtoList, &UserDTO{
			ID:     user.ID,
			Name:   user.Name,
			Rating: user.Rating,
		})
	}

	return dtoList, nil
}

func (a *AuthService) generateJWT(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(a.expirationTime).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(a.secretKey)
}
