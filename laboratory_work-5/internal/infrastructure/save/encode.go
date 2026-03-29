package save

import (
	"encoding/json"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

func newEncoder(format string) (Encoder, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		return JSONEncoder{}, nil
	case "yaml":
		return YAMLEncoder{}, nil
	default:
		return nil, ErrUnsupportedEncoder
	}
}

type JSONEncoder struct{}

func (j JSONEncoder) Encode(w io.Writer, data any) error {
	e := json.NewEncoder(w)
	e.SetIndent("", " ")

	return e.Encode(data)
}

type YAMLEncoder struct{}

func (y YAMLEncoder) Encode(w io.Writer, data any) error {
	e := yaml.NewEncoder(w)
	e.SetIndent(1)

	return e.Encode(data)
}
