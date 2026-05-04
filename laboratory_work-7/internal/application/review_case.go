package application

import (
	"context"

	"github.com/google/uuid"
)

type ReviewService struct{}

// check: userID and bookID must exist
func (r *ReviewService) AddBookReview(ctx context.Context, userID, bookID uuid.UUID, mark uint, report string) (BookReviewDTO, error) {
	return BookReviewDTO{}, nil
}

// check: fromUserID and toUserID must exist
func (r *ReviewService) AddUserReview(ctx context.Context, fromUserID, toUserID uuid.UUID, mark uint, report string) (UserReviewDTO, error) {
	return UserReviewDTO{}, nil
}
