package domain

import (
	"github.com/google/uuid"
)

type Review struct {
	id     uuid.UUID
	fromID uuid.UUID
	toID   uuid.UUID
	report string
	mark   uint
}

func NewReview(fromID, toID uuid.UUID, mark uint, report string) (*Review, error) {
	if fromID == toID {
		return nil, ErrSelfReview
	}

	if mark > 10 {
		return nil, NewErrMark(mark)
	}

	return &Review{
		id:     uuid.New(),
		fromID: fromID,
		toID:   toID,
		report: report,
		mark:   mark,
	}, nil
}

func (u *Review) ID() uuid.UUID {
	return u.id
}

func (u *Review) FromID() uuid.UUID {
	return u.fromID
}

func (u *Review) ToID() uuid.UUID {
	return u.toID
}

func (u *Review) Report() string {
	return u.report
}

func (u *Review) ChangeReport(newReport string) {
	u.report = newReport
}

func (u *Review) Mark() uint {
	return u.mark
}

func (u *Review) ChangeMark(newMark uint) error {
	if newMark > 10 {
		return NewErrMark(newMark)
	}

	u.mark = newMark
	return nil
}

type BookReview struct {
	r *Review
}

func NewBookReview(userID, bookID uuid.UUID, mark uint, report string) (*BookReview, error) {
	b := &BookReview{}

	r, err := NewReview(userID, bookID, mark, report)
	if err != nil {
		return nil, err
	}

	b.r = r
	return b, nil
}

func (b *BookReview) UserID() uuid.UUID {
	return b.r.fromID
}

func (b *BookReview) BookID() uuid.UUID {
	return b.r.toID
}

type UserReview struct {
	r *Review
}

func NewUserReview(fromUserID, toUserID uuid.UUID, mark uint, report string) (*BookReview, error) {
	b := &BookReview{}

	r, err := NewReview(fromUserID, toUserID, mark, report)
	if err != nil {
		return nil, err
	}

	b.r = r
	return b, nil
}

func (b *BookReview) FromID() uuid.UUID {
	return b.r.fromID
}

func (b *BookReview) ToID() uuid.UUID {
	return b.r.toID
}
