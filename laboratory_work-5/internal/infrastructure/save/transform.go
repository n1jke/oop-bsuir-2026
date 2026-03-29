package save

import (
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"strings"
)

func newTransformers(cfg ExportConfig) ([]Transformer, error) {
	if len(cfg.Transformations) == 0 {
		return nil, nil
	}

	out := make([]Transformer, 0, len(cfg.Transformations))
	for i := range cfg.Transformations {
		switch strings.ToLower(strings.TrimSpace(cfg.Transformations[i])) {
		case "encrypt":
			out = append(out, AESCipher{})
		case "compress":
			out = append(out, GzipCompress{})
		default:
			return nil, ErrUnsupportedTransform
		}
	}

	return out, nil
}

type AESCipher struct{}

func (AESCipher) Wrap(dst io.Writer) (io.WriteCloser, error) {
	// TODO: replace hardcoded key with key management provider.
	key, err := hex.DecodeString("01037d8224518c11bd2fd5cec6f1003530b7b49359dd71dcf931b5be36fcc305")
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// IV so output stream can be decrypted
	if _, err := dst.Write(iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)

	return nopWriteCloser{Writer: &cipher.StreamWriter{S: stream, W: dst}}, nil
}

type GzipCompress struct{}

func (GzipCompress) Wrap(dst io.Writer) (io.WriteCloser, error) {
	return gzip.NewWriter(dst), nil
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error {
	return nil
}
