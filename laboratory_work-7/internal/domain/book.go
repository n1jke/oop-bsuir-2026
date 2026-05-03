package domain

import (
	"slices"

	"github.com/google/uuid"
)

const MaxDescriptionLength = 255

type ISBN string

func NewISBN(ident string) (ISBN, error) {
	isbn := ISBN(ident)

	if len(ident) != 10 && len(ident) != 13 {
		return "", ErrInvalidISBN
	}

	return isbn, nil
}

type Book struct {
	id          uuid.UUID
	title       string
	authorID    uuid.UUID
	isbn        ISBN
	description string
	topics      []string
}

func NewBook(title string, authorID uuid.UUID, isbn ISBN, description string, topics ...string) (*Book, error) {
	if len(description) > MaxDescriptionLength {
		return nil, ErrLongDescription
	}

	return &Book{
		id:          uuid.New(),
		title:       title,
		authorID:    authorID,
		isbn:        isbn,
		description: description,
		topics:      topics,
	}, nil
}

func (b *Book) ID() uuid.UUID {
	return b.id
}

func (b *Book) Title() string {
	return b.title
}

func (b *Book) AuthorID() uuid.UUID {
	return b.authorID
}

func (b *Book) ISBN() ISBN {
	return b.isbn
}

func (b *Book) Description() string {
	return b.description
}

func (b *Book) UpdateDescription(descr string) error {
	if len(descr) > MaxDescriptionLength {
		return ErrLongDescription
	}

	b.description = descr
	return nil
}

func (b *Book) Topics() []string {
	return slices.Clone(b.topics)
}

func (b *Book) AddTopic(topic string) {
	b.topics = append(b.topics, topic)
}

func (b *Book) CleanTopics() {
	b.topics = nil
}

type OwnedBook struct {
	id      uuid.UUID
	bookID  uuid.UUID
	ownerID uuid.UUID
	status  OwnedBookStatus
}

type OwnedBookStatus int

const (
	Available OwnedBookStatus = iota
	Reserved
	Rent
	Hidden
)
