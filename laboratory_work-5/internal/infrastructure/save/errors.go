package save

import "errors"

var (
	ErrUnsupportedEncoder   = errors.New("unsupported encoder format")
	ErrUnsupportedTransform = errors.New("unsupported transformer")
	ErrEmptyOutputPath      = errors.New("output path must not be empty")
	ErrNilEncoder           = errors.New("encoder must not be nil")
	ErrInvalidPayload       = errors.New("response has no data to save")
)
