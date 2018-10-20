package qr

import (
	"bytes"
	"fmt"
	"github.com/boombuler/barcode"
	"image/png"
	"io"

	"github.com/boombuler/barcode/qr"
)

type writer struct {
	io.WriteCloser
	f   *Format
	w   io.Writer
	buf bytes.Buffer
}

// NewWriter returns a qr code format writer instance.
func NewWriter(w io.Writer, f *Format) io.WriteCloser {
	return &writer{w: w, f: f}
}

func (w *writer) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}

func (w *writer) Close() error {
	qrCode, err := qr.Encode(w.buf.String(), w.f.Level, qr.AlphaNumeric)
	if err != nil {
		return fmt.Errorf("failed to create qr code: %v", err)
	}

	qrCode, err = barcode.Scale(qrCode, w.f.Size, w.f.Size)
	if err != nil {
		return fmt.Errorf("failed to scale qr code: %v", err)
	}

	return png.Encode(w.w, qrCode)
}
