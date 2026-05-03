package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/domain"
)

func TestNewReview(t *testing.T) {
	selfReview := uuid.New()

	tests := []struct {
		name    string
		fromID  uuid.UUID
		toID    uuid.UUID
		mark    uint
		wantErr error
	}{
		{"valid", uuid.New(), uuid.New(), 8, nil},
		{"valid max mark", uuid.New(), uuid.New(), 10, nil},
		{"invalid selfreview", selfReview, selfReview, 5, domain.ErrSelfReview},
		{"invalid mark", uuid.New(), uuid.New(), 11, domain.NewErrMark(11)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			review, err := domain.NewReview(tt.fromID, tt.toID, tt.mark, "report")

			if tt.wantErr != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, review)
		})
	}
}

func TestReview_ChangeMark(t *testing.T) {
	tests := []struct {
		name    string
		initial uint
		new     uint
		wantErr bool
	}{
		{"valid change", 5, 8, false},
		{"valid to zero", 5, 0, false},
		{"invalid mark", 5, 11, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			review, err := domain.NewReview(uuid.New(), uuid.New(), tt.initial, "original")
			require.NoError(t, err)

			err = review.ChangeMark(tt.new)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.new, review.Mark())
		})
	}
}

func TestNewBookReview(t *testing.T) {
	t.Run("ValidBookReview_Success", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()
		bookID := uuid.New()

		review, err := domain.NewBookReview(userID, bookID, 9, "good book for system design prep")

		require.NoError(t, err)
		require.Equal(t, userID, review.UserID())
		require.Equal(t, bookID, review.BookID())
	})
}
