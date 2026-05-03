package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidISBN            = errors.New("invalid ISBN identifier")
	ErrLongDescription        = errors.New("description is too long")
	ErrSelfReview             = errors.New("user cannot review themselves")
	ErrInvalidStatus          = errors.New("invalid exchange status move")
	ErrDatePolicyNotConfigure = errors.New("DatePolicy not configure")
)

// ErrMark
type ErrMark struct {
	Mark uint
}

func NewErrMark(mark uint) error {
	return ErrMark{Mark: mark}
}

func (e ErrMark) Error() string {
	return fmt.Sprintf("invalid mark %d, should be in 0-10 range", e.Mark)
}
