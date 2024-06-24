package base64

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/corvus-ch/horcrux/input"
	"github.com/corvus-ch/horcrux/output"
)

// Name holds the name of the Format.
const Name = "base64"

// New creates a new Format of type base64.
func New(input input.Input) *Format {
	return &Format{input: input}
}

// Format represents the base64 type format.
type Format struct {
	input input.Input
}

// OutputFileName returns the file name for the given x.
func (f *Format) OutputFileName(x byte) string {
	return fmt.Sprintf("%s.base64.%03d", f.input.Stem(), x)
}

// Writer creates a new base64 format writer for the part identified by x.
func (f *Format) Writer(x byte, out output.Output) (io.Writer, []io.Closer, error) {
	path := f.OutputFileName(x)

	file, err := os.Create(path)
	if nil != err {
		return nil, nil, err
	}

	enc := base64.NewEncoder(base64.StdEncoding, file)

	close(out.Append(Name, path, nil))

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
