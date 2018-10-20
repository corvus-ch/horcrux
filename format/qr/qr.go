package qr

import (
	"fmt"
	"io"

	"github.com/boombuler/barcode/qr"
	"github.com/corvus-ch/zbase32"
)

// Name holds the name of the Format.
const Name = "qr"

const alphabet = "YBNDRFG8EJKMCPQXOT1UWISZA345H769"

var encoding = zbase32.NewEncoding(alphabet)

// New creates a new Format of type raw.
func New(stem string) *Format {
	return &Format{
		stem:  stem,
		Level: qr.M,
		Size:  500,
	}
}

// Format represents the raw type format.
type Format struct {
	stem  string
	Level qr.ErrorCorrectionLevel
	Size  int
}

// OutputFileName returns the file name for the given x.
func (f *Format) OutputFileName(x byte) string {
	return fmt.Sprintf("%s.%03d.png", f.stem, x)
}

// Writer creates a new Format writer using the given writer as output.
func (f *Format) Writer(out io.Writer) (io.Writer, []io.Closer, error) {
	cs := make([]io.Closer, 2)
	w := NewWriter(out, f)
	cs[0] = w
	enc := zbase32.NewEncoder(encoding, w)
	cs[1] = enc

	return enc, cs, nil
}

// Reader creates a new Format reader using the given reader as input.
func (f *Format) Reader(r io.Reader) (io.Reader, error) {
	return nil, fmt.Errorf("format does not yet have read support")
}

// Name returns the formats name.
func (f *Format) Name() string {
	return Name
}
