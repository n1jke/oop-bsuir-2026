package save

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/application"
)

type ExportConfig struct {
	Format          string
	Transformations []string
	OutPath         string
}

type ResponseSaver struct {
	encoder      Encoder
	transformers []Transformer
	outPath      string
}

type Encoder interface {
	Encode(w io.Writer, v any) error
}

type Transformer interface {
	Wrap(dst io.Writer) (io.WriteCloser, error)
}

func NewResponseSaver(cfg ExportConfig) (*ResponseSaver, error) {
	encoder, err := newEncoder(cfg.Format)
	if err != nil {
		return nil, err
	}

	transformers, err := newTransformers(cfg)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(cfg.OutPath) == "" {
		return nil, ErrEmptyOutputPath
	}

	return &ResponseSaver{
		encoder:      encoder,
		transformers: transformers,
		outPath:      cfg.OutPath,
	}, nil
}

func (rs *ResponseSaver) Save(data *application.ServiceResponse) error {
	if data == nil {
		return ErrInvalidPayload
	}

	file, err := os.Create(rs.outPath)
	if err != nil {
		return err
	}

	defer func() {
		inErr := file.Close()
		err = errors.Join(err, inErr)
	}()

	writer := io.Writer(file)
	closers := make([]io.Closer, 0, len(rs.transformers))

	// reverse traversal to kepp T1(T2(...Tn(data)))
	for i := len(rs.transformers) - 1; i >= 0; i-- {
		wrapped, err := rs.transformers[i].Wrap(writer)
		if err != nil {
			return err
		}

		closers = append(closers, wrapped)
		writer = wrapped
	}

	// traversal encode(T1(...Tn(data)))
	encodeErr := rs.encoder.Encode(writer, mapServiceResponse(data))
	closeErr := closeWriters(closers)

	if encodeErr != nil {
		return errors.Join(err, encodeErr)
	}

	if closeErr != nil {
		return errors.Join(err, closeErr)
	}

	return nil
}

func closeWriters(closers []io.Closer) error {
	var rErr error

	for i := len(closers) - 1; i >= 0; i-- {
		if err := closers[i].Close(); err != nil {
			rErr = errors.Join(rErr, err)
		}
	}

	return rErr
}
