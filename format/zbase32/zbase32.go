package zbase32

import (
	"fmt"
	"io"
	"os"

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

// Writer creates a new zbase32 format writer for the part identified by x.
func (f *Format) Writer(x byte) (io.Writer, []io.Closer, error) {
	file, err := os.Create(f.OutputFileName(x))
	if nil != err {
		return nil, nil, err
	}

	enc := zbase32.NewEncoder(zbase32.StdEncoding, file)

	return enc, []io.Closer{file, enc}, nil
}

// Reader creates a new Format reader using the given reader as input.
func (f *Format) Reader(r io.Reader) (io.Reader, error) {
	return zbase32.NewDecoder(zbase32.StdEncoding, r), nil
}

// Name returns the formats name.
func (f *Format) Name() string {
	return Name
}
