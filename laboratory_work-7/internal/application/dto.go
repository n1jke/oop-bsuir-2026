package application

import (
	"time"

	"github.com/google/uuid"
)

// api dto

type UserDTO struct {
	ID     uuid.UUID
	Name   string
	Rating float64
}

type BookDTO struct {
	ID          uuid.UUID
	Title       string
	AuthorName  string
	ISBN        string
	Description string
	Topics      []string
}

type BookListResponse struct {
	Books []BookDTO
	Total int
}

type ExchangeDTO struct {
	ID          uuid.UUID
	OwnedBookID uuid.UUID
	FromID      uuid.UUID
	ToID        uuid.UUID
	Status      string
	CreatedAt   time.Time
	ExpiresAt   time.Time
	Note        string
}

type BookReviewDTO struct {
	ID     uuid.UUID
	FromID uuid.UUID
	BookID uuid.UUID
	Mark   uint
	Report string
}

type UserReviewDTO struct {
	ID     uuid.UUID
	FromID uuid.UUID
	ToID   uuid.UUID
	Mark   uint
	Report string
}

type LoginResponse struct {
	Token string
	User  UserDTO
}

type RegisterResponse struct {
	ID    uuid.UUID
	Token string
}

// infra dto

type UserRepoDTO struct {
	ID           uuid.UUID
	Name         string
	PasswordHash string
	Rating       float64
}
