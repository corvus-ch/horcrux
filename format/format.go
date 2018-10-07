package format

import (
	"fmt"
	"io"
)

var DEFAULT = "raw"

// Format describes the interface for the various input/output formats.
type Format interface {
	// OutputFileName returns the file name for the given x.
	OutputFileName(x byte) string

	// Writer creates a new format writer using the given writer as output.
	Writer(w io.Writer) (io.Writer, []io.Closer, error)

	// Reader creates a new format reader using the given reader as input.
	Reader(r io.Reader) (io.Reader, error)
}

func New(format, stem string) (Format, error) {
	switch format {

	// TODO

	default:
		return nil, fmt.Errorf("unknown format %s", format)
	}
}
