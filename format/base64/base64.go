package base64

import (
	"encoding/base64"
	"fmt"
	"io"
)

// Name holds the name of the Format.
const Name = "base64"

// New creates a new Format of type base64.
func New(stem string) *Format {
	return &Format{Stem: stem}
}

// Format represents the base64 type format.
type Format struct {
	Stem string
}

// OutputFileName returns the file name for the given x.
func (f *Format) OutputFileName(x byte) string {
	return fmt.Sprintf("%s.base64.%03d", f.Stem, x)
}

// Writer creates a new Format writer using the given writer as output.
func (f *Format) Writer(out io.Writer) (io.Writer, []io.Closer, error) {
	enc := base64.NewEncoder(base64.StdEncoding, out)

	return enc, []io.Closer{enc}, nil
}

// Reader creates a new Format reader using the given reader as input.
func (f *Format) Reader(r io.Reader) (io.Reader, error) {
	return base64.NewDecoder(base64.StdEncoding, r), nil
}

// Name returns the formats name.
func (f *Format) Name() string {
	return Name
}
