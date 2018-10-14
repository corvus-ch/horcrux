package text

import (
	"fmt"
	"io"

	"gopkg.in/corvus-ch/zbase32.v1"
)

// Name holds the name of the format.
const Name = "text"

// New creates a new format of type Text.
func New(stem string) *Format {
	return &Format{Stem: stem, LineLength: 72}
}

// Format represents the text type format.
type Format struct {
	LineLength uint8
	Stem       string
}

// OutputFileName returns the file name for the given x.
func (f *Format) OutputFileName(x byte) string {
	return fmt.Sprintf("%s.txt.%03d", f.Stem, x)
}

// Writer creates a new format writer using the given writer as output.
func (f *Format) Writer(out io.Writer) (io.Writer, []io.Closer, error) {
	cs := make([]io.Closer, 2)
	w := NewWriter(out, f)
	cs[0] = w
	enc := zbase32.NewEncoder(zbase32.StdEncoding, w)
	cs[1] = enc

	return enc, cs, nil
}

// Reader creates a new format reader using the given reader as input.
func (f *Format) Reader(r io.Reader) (io.Reader, error) {
	return zbase32.NewDecoder(zbase32.StdEncoding, NewReader(r)), nil
}

// Name returns the formats name.
func (f *Format) Name() string {
	return Name
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
