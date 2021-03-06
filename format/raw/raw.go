package raw

import (
	"fmt"
	"io"
	"os"

	"github.com/corvus-ch/horcrux/input"
)

// Name holds the name of the Format.
const Name = "raw"

// New creates a new Format of type raw.
func New(input input.Input) *Format {
	return &Format{input: input}
}

// Format represents the raw type format.
type Format struct {
	input input.Input
}

// OutputFileName returns the file name for the given x.
func (f *Format) OutputFileName(x byte) string {
	return fmt.Sprintf("%s.raw.%03d", f.input.Stem(), x)
}

// Writer creates a new raw format writer for the part identified by x.
func (f *Format) Writer(x byte) (io.Writer, []io.Closer, error) {
	file, err := os.Create(f.OutputFileName(x))
	if nil != err {
		return nil, nil, err
	}

	return file, []io.Closer{file}, nil
}

// Reader creates a new Format reader using the given reader as input.
func (f *Format) Reader(r io.Reader) (io.Reader, error) {
	return r, nil
}

// Name returns the formats name.
func (f *Format) Name() string {
	return Name
}
