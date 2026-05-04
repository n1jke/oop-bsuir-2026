package application

import (
	"time"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID     uuid.UUID
	Name   string
	Rating float64
}

type BookDTO struct {
	Id          uuid.UUID
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
	Id          uuid.UUID
	OwnedBookId uuid.UUID
	FromId      uuid.UUID
	ToId        uuid.UUID
	Status      string
	CreatedAt   time.Time
	ExpiresAt   time.Time
	Note        string
}

type ReviewDTO struct {
	Id     uuid.UUID
	FromId uuid.UUID
	ToId   uuid.UUID
	BookId uuid.UUID
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
