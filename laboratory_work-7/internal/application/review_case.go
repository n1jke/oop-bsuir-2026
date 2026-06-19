package application

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/domain"
)

type ReviewService struct {
	logger     *slog.Logger
	reviewRepo ReviewRepository
	bookRepo   BookRepository
	userRepo   UserRepository
	tx         Transactor
}

func NewReviewService(logger *slog.Logger, reviewRepo ReviewRepository, bookRepo BookRepository, userRepo UserRepository,
	tx Transactor,
) *ReviewService {
	logger = logger.With("module", "review-service")

	return &ReviewService{
		logger:     logger,
		reviewRepo: reviewRepo,
		bookRepo:   bookRepo,
		userRepo:   userRepo,
		tx:         tx,
	}
}

func (r *ReviewService) AddBookReview(ctx context.Context, userID, bookID uuid.UUID, mark uint, rep string) (*BookReviewDTO, error) {
	var (
		review *domain.BookReview
		errIn  error
	)

	err := r.tx.WithTransaction(ctx, func(context.Context) error {
		review, errIn = r.addBookReviewTx(ctx, userID, bookID, mark, rep)
		return errIn
	})
	if err != nil {
		return nil, err
	}

	return &BookReviewDTO{
		ID:     review.ID(),
		FromID: review.UserID(),
		BookID: review.BookID(),
		Mark:   review.Mark(),
		Report: review.Report(),
	}, nil
}

func (r *ReviewService) AddUserReview(ctx context.Context, fromUserID, toUserID uuid.UUID, mark uint, rep string) (*UserReviewDTO, error) {
	var (
		review *domain.UserReview
		errIn  error
	)

	err := r.tx.WithTransaction(ctx, func(context.Context) error {
		review, errIn = r.addUserReviewTx(ctx, fromUserID, toUserID, mark, rep)
		return errIn
	})
	if err != nil {
		return nil, err
	}

	return &UserReviewDTO{
		ID:     review.ID(),
		FromID: review.FromID(),
		ToID:   review.ToID(),
		Mark:   review.Mark(),
		Report: review.Report(),
	}, nil
}
