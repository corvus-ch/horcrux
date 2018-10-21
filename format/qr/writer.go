package qr

import (
	"bytes"
	"fmt"
	"github.com/boombuler/barcode"
	"image/png"
	"io"
	"os"

	"github.com/boombuler/barcode/qr"
)

type writer struct {
	io.WriteCloser
	f   *Format
	buf bytes.Buffer
	x   byte
	n   int
}

// NewWriter returns a qr code format writer instance.
func NewWriter(f *Format, x byte) io.WriteCloser {
	return &writer{f: f, x: x}
}

func (w *writer) Write(p []byte) (int, error) {
	n, err := w.buf.Write(p)
	if err != nil {
		return n, err
	}

	for w.buf.Len() >= 3391 {
		if err := w.createImage(); err != nil {
			return n, err
		}
	}

	return n, nil
}

func (w *writer) Close() error {
	if w.buf.Len() > 0 {
		return w.createImage()
	}

	return nil
}

func (w *writer) createImage() error {
	data := w.buf.Next(3391)
	code, err := qr.Encode(string(data), w.f.Level, qr.AlphaNumeric)
	if err != nil {
		return fmt.Errorf("failed to create qr code: %v", err)
	}

	code, err = barcode.Scale(code, w.f.Size, w.f.Size)
	if err != nil {
		return fmt.Errorf("failed to scale qr code: %v", err)
	}

	// create the output file
	w.n++
	file, err := os.Create(fmt.Sprintf("%s.%03d.%d.png", w.f.Stem, w.x, w.n))
	if err != nil {
		return fmt.Errorf("failed to open output file: %v", err)
	}
	defer file.Close()

	return png.Encode(file, code)
}
