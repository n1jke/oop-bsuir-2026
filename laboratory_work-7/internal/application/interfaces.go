package application

import (
	"context"
)

//go:generate mockgen -source=interfaces.go -destination=mock/mock.go -package=mocks

type Transactor interface {
	WithTransaction(context.Context, func(context.Context) error) error
}
