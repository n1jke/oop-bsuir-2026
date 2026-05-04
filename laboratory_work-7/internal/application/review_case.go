package application

import (
	"context"

	"github.com/google/uuid"
)

type ReviewService struct{}

func (r *ReviewService) AddBookReview(ctx context.Context, userID, bookID uuid.UUID, mark uint, report string) (ReviewDTO, error) {
	return ReviewDTO{}, nil
}

func (r *ReviewService) AddUserReview(ctx context.Context, fromUserID, toUserID uuid.UUID, mark uint, report string) (ReviewDTO, error) {
	return ReviewDTO{}, nil
}
