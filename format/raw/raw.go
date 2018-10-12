package raw

import (
	"fmt"
	"io"
)

// Name holds the name of the format.
const Name = "raw"

// New creates a new format of type raw.
func New(stem string) *format {
	return &format{stem: stem}
}

type format struct {
	stem string
}

// OutputFileName returns the file name for the given x.
func (f *format) OutputFileName(x byte) string {
	return fmt.Sprintf("%s.raw.%03d", f.stem, x)
}

// Writer creates a new format writer using the given writer as output.
func (f *format) Writer(out io.Writer) (io.Writer, []io.Closer, error) {
	return out, nil, nil
}

// Reader creates a new format reader using the given reader as input.
func (f *format) Reader(r io.Reader) (io.Reader, error) {
	return r, nil
}

// Name returns the formats name.
func (f *format) Name() string {
	return Name
}
