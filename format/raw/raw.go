package raw

import (
	"fmt"
	"io"
)

const Name = "raw"

func New(stem string) *format {
	return &format{stem: stem}
}

type format struct {
	stem string
}

func (f *format) OutputFileName(x byte) string {
	return fmt.Sprintf("%s.raw.%03d", f.stem, x)
}

func (f *format) Writer(out io.Writer) (io.Writer, []io.Closer, error) {
	return out, nil, nil
}

func (f *format) Reader(r io.Reader) (io.Reader, error) {
	return r, nil
}

func (f *format) Name() string {
	return Name
}
