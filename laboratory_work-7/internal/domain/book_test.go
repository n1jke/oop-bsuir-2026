package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/domain"
)

func TestNewISBN(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid digits=10", "0-306-40615-5", false},
		{"valid digits=13", "978-1-60309-542-6", false},
		{"indalid digits=8", "12-3-45678", true},
		{"invalid digits=14", "978-1-60309-54216", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			isbn, err := domain.NewISBN(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, domain.ISBN(tt.input), isbn)
		})
	}
}

func TestNewBook(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		isbn        domain.ISBN
		description string
		wantErr     bool
	}{
		{"valid", "Go Programming Language", "1234567890", "guide", false},
		{"valid with max description", "Book Title", "1234567890123", string(make([]byte, domain.MaxDescriptionLength)), false},
		{"invalid long description", "Book Title", "1234567890", string(make([]byte, domain.MaxDescriptionLength+1)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			book, err := domain.NewBook(tt.title, uuid.New(), tt.isbn, tt.description, "topic1", "topic2")

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.title, book.Title())
			require.Equal(t, tt.description, book.Description())
		})
	}
}

func TestBook_UpdateDescription(t *testing.T) {
	tests := []struct {
		name    string
		newDesc string
		wantErr bool
	}{
		{"valid empty description", "", false},
		{"valid", string(make([]byte, domain.MaxDescriptionLength)), false},
		{"invalid long description", string(make([]byte, domain.MaxDescriptionLength+1)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			book, err := domain.NewBook("DDD", uuid.New(), "1234567890", "init")
			require.NoError(t, err)

			err = book.UpdateDescription(tt.newDesc)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.newDesc, book.Description())
		})
	}
}

func TestBook_Topics(t *testing.T) {
	t.Run("immutable topics", func(t *testing.T) {
		t.Parallel()

		book, err := domain.NewBook("Clean architecture", uuid.New(), "1234567890", "architecture", "swe", "learn")
		require.NoError(t, err)

		topics := book.Topics()
		topics = append(topics, "dumped")

		require.Len(t, book.Topics(), 2)
		require.Len(t, topics, 3)
	})
}
