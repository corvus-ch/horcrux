package text

import (
	"fmt"
	"io"
	"os"

	"github.com/corvus-ch/horcrux/input"
	"gopkg.in/corvus-ch/zbase32.v1"
)

// Name holds the name of the format.
const Name = "text"

// New creates a new format of type Text.
func New(input input.Input) *Format {
	return &Format{input: input, LineLength: 72}
}

// Format represents the text type format.
type Format struct {
	input      input.Input
	LineLength uint8
}

// OutputFileName returns the file name for the given x.
func (f *Format) OutputFileName(x byte) string {
	return fmt.Sprintf("%s.txt.%03d", f.input.Stem(), x)
}

// Writer creates a new text format writer for the part identified by x.
func (f *Format) Writer(x byte) (io.Writer, []io.Closer, error) {
	file, err := os.Create(f.OutputFileName(x))
	if nil != err {
		return nil, nil, err
	}

	w, err := NewWriter(file, f, NewData(f.input, x))
	if err != nil {
		return nil, []io.Closer{file}, err
	}

	enc := zbase32.NewEncoder(zbase32.StdEncoding, w)

	return enc, []io.Closer{file, w, enc}, nil
}

// Reader creates a new format reader using the given reader as input.
func (f *Format) Reader(r io.Reader) (io.Reader, error) {
	return zbase32.NewDecoder(zbase32.StdEncoding, NewReader(r)), nil
}

// Name returns the formats name.
func (f *Format) Name() string {
	return Name
}
