package zbase32

import (
	"fmt"
	"io"

	"gopkg.in/corvus-ch/zbase32.v1"
)

// Name holds the name of the Format.
const Name = "zbase32"

// New creates a new Format of type zbase32.
func New(stem string) *Format {
	return &Format{Stem: stem}
}

// Format represents the zbase32 type format.
type Format struct {
	Stem string
}

// OutputFileName returns the file name for the given x.
func (f *Format) OutputFileName(x byte) string {
	return fmt.Sprintf("%s.zbase32.%03d", f.Stem, x)
}

// Writer creates a new Format writer using the given writer as output.
func (f *Format) Writer(out io.Writer) (io.Writer, []io.Closer, error) {
	enc := zbase32.NewEncoder(zbase32.StdEncoding, out)

	return enc, []io.Closer{enc}, nil
}

// Reader creates a new Format reader using the given reader as input.
func (f *Format) Reader(r io.Reader) (io.Reader, error) {
	return zbase32.NewDecoder(zbase32.StdEncoding, r), nil
}

// Name returns the formats name.
func (f *Format) Name() string {
	return Name
}
