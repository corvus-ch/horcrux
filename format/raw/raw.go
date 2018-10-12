package raw

import (
	"fmt"
	"io"
)

// Name holds the name of the Format.
const Name = "raw"

// New creates a new Format of type raw.
func New(stem string) *Format {
	return &Format{stem: stem}
}

type Format struct {
	stem string
}

// OutputFileName returns the file name for the given x.
func (f *Format) OutputFileName(x byte) string {
	return fmt.Sprintf("%s.raw.%03d", f.stem, x)
}

// Writer creates a new Format writer using the given writer as output.
func (f *Format) Writer(out io.Writer) (io.Writer, []io.Closer, error) {
	return out, nil, nil
}

// Reader creates a new Format reader using the given reader as input.
func (f *Format) Reader(r io.Reader) (io.Reader, error) {
	return r, nil
}

// Name returns the formats name.
func (f *Format) Name() string {
	return Name
}
