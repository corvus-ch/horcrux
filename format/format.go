package format

import (
	"fmt"
	"github.com/corvus-ch/horcrux/format/qr"
	"github.com/corvus-ch/horcrux/format/text"
	"io"

	"github.com/corvus-ch/horcrux/format/base64"
	"github.com/corvus-ch/horcrux/format/raw"
	"github.com/corvus-ch/horcrux/format/zbase32"
)

// Default holds the name of the default format.
const Default = text.Name

// Format describes the interface for the various input/output formats.
type Format interface {
	// OutputFileName returns the file name for the given x.
	OutputFileName(x byte) string

	// Writer creates a new format writer using the given writer as output.
	Writer(w io.Writer) (io.Writer, []io.Closer, error)

	// Reader creates a new format reader using the given reader as input.
	Reader(r io.Reader) (io.Reader, error)

	// Name returns the formats name.
	Name() string
}

// New creates a new format object.
func New(format, stem string) (Format, error) {
	switch format {
	case raw.Name:
		return raw.New(stem), nil

	case zbase32.Name:
		return zbase32.New(stem), nil

	case base64.Name:
		return base64.New(stem), nil

	case qr.Name:
		return qr.New(stem), nil

	case text.Name:
		return text.New(stem), nil

	default:
		return nil, fmt.Errorf("unknown format %s", format)
	}
}
