package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/domain"
)

func (r *ReviewService) addBookReviewTx(ctx context.Context, userID, bookID uuid.UUID, mark uint, rep string) (*domain.BookReview, error) {
	_, err := r.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrUserNotFound
		}

		r.logger.Error("get by id from userRepo", slog.Any("err", err))

		return nil, ErrUnvailible
	}

	_, err = r.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrBookNotFound
		}

		r.logger.Error("get by id from bookRepo", slog.Any("err", err))

		return nil, ErrUnvailible
	}

	review, err := domain.NewBookReview(userID, bookID, mark, rep)
	if err != nil {
		return nil, fmt.Errorf("invalid params: %w", err)
	}

	err = r.reviewRepo.AddBookReview(ctx, *review)
	if err != nil {
		r.logger.Error("add book review to repo", slog.Any("reviewID", review.ID()), slog.Any("err", err))
		return nil, ErrUnvailible
	}

	return review, nil
}

func (r *ReviewService) addUserReviewTx(ctx context.Context, fromID, toID uuid.UUID, mark uint, rep string) (*domain.UserReview, error) {
	_, err := r.userRepo.GetByID(ctx, fromID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrUserNotFound
		}

		r.logger.Error("get by id from userRepo", slog.Any("err", err))

		return nil, ErrUnvailible
	}

	_, err = r.userRepo.GetByID(ctx, toID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrUserNotFound
		}

		r.logger.Error("get by id from userRepo", slog.Any("err", err))

		return nil, ErrUnvailible
	}

	review, err := domain.NewUserReview(fromID, toID, mark, rep)
	if err != nil {
		return nil, fmt.Errorf("invalid params: %w", err)
	}

	err = r.reviewRepo.AddUserReview(ctx, *review)
	if err != nil {
		r.logger.Error("add user review to repo", slog.Any("reviewID", review.ID()), slog.Any("err", err))
		return nil, ErrUnvailible
	}

	return review, nil
}
