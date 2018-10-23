package qr

import (
	"fmt"
	"io"

	"github.com/corvus-ch/horcrux/meta"
	"github.com/corvus-ch/zbase32"
)

// Name holds the name of the Format.
const Name = "qr"

const alphabet = "YBNDRFG8EJKMCPQXOT1UWISZA345H769"

var encoding = zbase32.NewEncoding(alphabet)

// New creates a new Format of type raw.
func New(input meta.Input) *Format {
	return &Format{
		in: input,
	}
}

// Format represents the raw type format.
type Format struct {
	in meta.Input
}

// OutputFileName returns the file name for the given x.
func (f *Format) OutputFileName(x byte) string {
	return fmt.Sprintf("%s.png.%03d", f.in.Stem(), x)
}

// Writer creates a new QR code format writer for the part identified by x.
func (f *Format) Writer(x byte) (io.Writer, []io.Closer, error) {
	w := NewWriter(f.in, x)
	enc := zbase32.NewEncoder(encoding, w)

	return enc, []io.Closer{w, enc}, nil
}

// Reader creates a new Format reader using the given reader as input.
func (f *Format) Reader(r io.Reader) (io.Reader, error) {
	return nil, fmt.Errorf("format does not yet have read support")
}

// Name returns the formats name.
func (f *Format) Name() string {
	return Name
}
