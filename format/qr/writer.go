package qr

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"math"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/corvus-ch/horcrux/meta"
)

type writer struct {
	io.WriteCloser
	buf   bytes.Buffer
	chunk int
	in    meta.Input
	level qr.ErrorCorrectionLevel
	n     int
	x     byte
}

// NewWriter returns a qr code format writer instance.
func NewWriter(in meta.Input, x byte) io.WriteCloser {
	return &writer{in: in, x: x, level: qr.M}
}

func (w *writer) Write(p []byte) (int, error) {
	n, err := w.buf.Write(p)
	if err != nil {
		return n, err
	}

	for w.buf.Len() >= w.ChunkSize() {
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
	data := w.buf.Next(w.ChunkSize())
	code, err := qr.Encode(string(data), w.level, qr.AlphaNumeric)
	if err != nil {
		return fmt.Errorf("failed to create qr code: %v", err)
	}

	code, err = barcode.Scale(code, 500, 500)
	if err != nil {
		return fmt.Errorf("failed to scale qr code: %v", err)
	}

	// create the output file
	w.n++
	file, err := os.Create(fmt.Sprintf("%s.%03d.%d.png", w.in.Stem(), w.x, w.n))
	if err != nil {
		return fmt.Errorf("failed to open output file: %v", err)
	}
	defer file.Close()

	return png.Encode(file, code)
}

// ChunkSize returns the number of encoded bytes written to a single qr code image.
func (w *writer) ChunkSize() int {
	if w.chunk != 0 {
		return w.chunk
	}

	capacity := w.Capacity()
	size := float64(w.in.Size() * 8 / 5)
	chunks := math.Ceil(size / float64(capacity))
	w.chunk = int(math.Ceil(size / chunks))
	if w.chunk > capacity || w.chunk < 0 {
		w.chunk = capacity
	}

	return w.chunk
}

// Capacity returns the number of bytes which fit into a single qr code image.
func (w *writer) Capacity() int {
	switch w.level {
	case qr.L:
		return 4296
	case qr.M:
		return 3391
	case qr.Q:
		return 2420
	default:
		return 1852
	}
}
