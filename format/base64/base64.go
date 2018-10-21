package base64

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
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

// Writer creates a new base64 format writer for the part identified by x.
func (f *Format) Writer(x byte) (io.Writer, []io.Closer, error) {
	file, err := os.Create(f.OutputFileName(x))
	if nil != err {
		return nil, nil, err
	}

	enc := base64.NewEncoder(base64.StdEncoding, file)

	return enc, []io.Closer{file, enc}, nil
}

// Reader creates a new Format reader using the given reader as input.
func (f *Format) Reader(r io.Reader) (io.Reader, error) {
	return base64.NewDecoder(base64.StdEncoding, r), nil
}

// Name returns the formats name.
func (f *Format) Name() string {
	return Name
}
