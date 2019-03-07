package qr

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"math"
	"os"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/corvus-ch/horcrux/input"
)

type writer struct {
	io.WriteCloser
	buf   bytes.Buffer
	in    input.Input
	level qr.ErrorCorrectionLevel
	n     int
	x     byte
}

// NewWriter returns a qr code format writer instance.
func NewWriter(in input.Input, x byte) io.WriteCloser {
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
	var data strings.Builder
	data.WriteString(w.index())
	data.Write(w.buf.Next(w.ChunkSize()))
	code, err := qr.Encode(data.String(), w.level, qr.AlphaNumeric)
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
	return ChunkSize(w.Capacity(), w.in.Size())
}

// Capacity returns the number of bytes which fit into a single qr code image.
func (w *writer) Capacity() int {
	indexLength := len(w.index())
	switch w.level {
	case qr.L:
		return 4296 - indexLength
	case qr.M:
		return 3391 - indexLength
	case qr.Q:
		return 2420 - indexLength
	default:
		return 1852 - indexLength
	}
}

func (w *writer) index() string {
	return fmt.Sprintf("%03d:%d::", w.x, w.n)
}

// NumChunks returns the number of images required to encode the data.
func NumChunks(capacity int, size int64) int {
	return int(math.Ceil(float64(size*8) / 5 / float64(capacity)))
}

// ChunkSize returns the number of bytes fitting into a single qr code image.
func ChunkSize(capacity int, size int64) int {
	chunks := NumChunks(capacity, size)
	chunk := int(math.Ceil(float64(size*8) / 5 / float64(chunks)))
	if chunk > capacity || chunk < 0 {
		chunk = capacity
	}

	return chunk
}
