package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/domain"
)

const eps = 0.001

func TestRating_GetRating(t *testing.T) {
	tests := []struct {
		name       string
		points     []uint
		wantRating float64
	}{
		{"0 to 0", []uint{}, 0.0},
		{"single 10", []uint{10}, 10.0},
		{"float num", []uint{5, 10}, 7.5},
		{"mixed scores", []uint{3, 7, 5, 9}, 6.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rating := &domain.Rating{}
			for _, p := range tt.points {
				rating.AddReview(p)
			}

			got := rating.GetRating()

			require.InDelta(t, tt.wantRating, got, eps)
		})
	}
}
