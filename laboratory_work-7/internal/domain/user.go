package domain

import "github.com/google/uuid"

type User struct {
	id     uuid.UUID
	name   string
	rating Rating
}

func NewUser(name string) *User {
	return &User{
		id:     uuid.New(),
		name:   name,
		rating: Rating{},
	}
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) AddReview(points uint) {
	u.rating.AddReview(points)
}

func (u *User) GetRating() float64 {
	return u.rating.GetRating()
}
