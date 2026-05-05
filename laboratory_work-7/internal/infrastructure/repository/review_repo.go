package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/domain"
)

type ReviewRepoSQL struct {
	db *pgxpool.Pool
}

const (
	insertBookReview = `
		INSERT INTO book_reviews(id, user_id, book_id, mark, report)
		VALUES ($1, $2, $3, $4, $5)
	`

	insertUserReview = `
		INSERT INTO user_reviews(id, from_id, to_id, mark, report)
		VALUES ($1, $2, $3, $4, $5)
	`

	selectBookReviews = `
		SELECT id, user_id, book_id, mark, report FROM book_reviews
		WHERE book_id = $1
	`

	selectUserReviews = `
		SELECT id, from_id, to_id, mark, report FROM user_reviews
		WHERE to_id = $1
	`
)

func (r ReviewRepoSQL) AddBookReview(ctx context.Context, review domain.BookReview) error {
	q := GetQuerier(ctx, r.db)

	_, err := q.Exec(ctx, insertBookReview,
		review.ID(), review.UserID(), review.BookID(), review.Mark(), review.Report(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r ReviewRepoSQL) AddUserReview(ctx context.Context, review domain.UserReview) error {
	q := GetQuerier(ctx, r.db)

	_, err := q.Exec(ctx, insertUserReview,
		review.ID(), review.FromID(), review.ToID(), review.Mark(), review.Report(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r ReviewRepoSQL) GetBookReviews(ctx context.Context, bookID uuid.UUID) ([]*domain.BookReview, error) {
	q := GetQuerier(ctx, r.db)

	rows, err := q.Query(ctx, selectBookReviews, bookID)
	if err != nil {
		return nil, err
	}

	reviews := make([]*domain.BookReview, 0)

	for rows.Next() {
		var (
			id, userID, bID uuid.UUID
			mark uint
			report string
		)

		if err := rows.Scan(&id, &userID, &bID, &mark, &report); err != nil {
			return nil, err
		}

		review, err := domain.CreateBookReview(id, userID, bID, mark, report)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r ReviewRepoSQL) GetUserReviews(ctx context.Context, userID uuid.UUID) ([]*domain.UserReview, error) {
	q := GetQuerier(ctx, r.db)

	rows, err := q.Query(ctx, selectUserReviews, userID)
	if err != nil {
		return nil, err
	}

	reviews := make([]*domain.UserReview, 0)

	for rows.Next() {
		var (
			id, fromID, toID uuid.UUID
			mark uint
			report string
		)

		if err := rows.Scan(&id, &fromID, &toID, &mark, &report); err != nil {
			return nil, err
		}

		review, err := domain.CreateUserReview(id, fromID, toID, mark, report)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func NewReviewRepoSQL(db *pgxpool.Pool) *ReviewRepoSQL {
	return &ReviewRepoSQL{db: db}
}
