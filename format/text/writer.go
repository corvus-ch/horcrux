package text

import (
	"fmt"
	"io"

	"github.com/martinlindhe/crc24"
)

type writer struct {
	io.WriteCloser
	err  error
	w    io.Writer    // The underlying writer
	n    int          // Lines count
	buf  []byte       // buffered data waiting to be encoded
	nbuf int          // number of bytes in buf
	crc  crc24.Hash24 // The checksum
}

// NewWriter returns an text format writer instance.
func NewWriter(w io.Writer, f *Format) io.WriteCloser {
	return &writer{
		w:   w,
		buf: make([]byte, bufLen(f.LineLength)),
		crc: crc24.New(),
	}
}

func (w *writer) Write(p []byte) (n int, err error) {
	if w.err != nil {
		return 0, w.err
	}

	// Leading fringe.
	if len(w.buf) > w.nbuf {
		var i int
		for i = 0; i < len(p) && w.nbuf < len(w.buf); i++ {
			w.buf[w.nbuf] = p[i]
			w.nbuf++
		}
		n += i
		p = p[i:]
		if w.nbuf < len(w.buf) {
			return
		}
		w.crc.Write(w.buf)
		if err := w.writeLine(w.buf); nil != err {
			return n, w.err
		}
		w.nbuf = 0
	}

	// Large interior chunks.
	for len(p) > len(w.buf) {
		nn := len(w.buf)
		w.crc.Write(p[:nn])
		if err := w.writeLine(p[:nn]); nil != err {
			return n, w.err
		}
		n += nn
		p = p[nn:]
	}

	// Trailing fringe.
	for i := 0; i < len(p); i++ {
		w.buf[i] = p[i]
	}
	w.nbuf = len(p)
	n += len(p)
	return
}

func (w *writer) Close() error {
	if w.nbuf > 0 {
		w.crc.Write(w.buf[:w.nbuf])
		if err := w.writeLine(w.buf[:w.nbuf]); nil != err {
			return err
		}
	}
	w.n++
	if _, err := w.w.Write([]byte(fmt.Sprintf("% 4d: %v", w.n, formatCRC24(w.crc)))); nil != err {
		return fmt.Errorf("Failed to write checksum: %v", err)
	}
	w.crc.Reset()
	return nil
}

func (w *writer) writeLine(data []byte) error {
	w.n++
	if _, err := w.w.Write([]byte(fmt.Sprintf("% 4d: ", w.n))); nil != err {
		return fmt.Errorf("Failed to write line index: %v", err)
	}
	for i := 0; i < len(data); i += 5 {
		j := min(i+5, len(data))
		if _, err := w.w.Write(data[i:j]); nil != err {
			return fmt.Errorf("Failed to write line data: %v", err)
		}
		if _, err := w.w.Write([]byte{' '}); nil != err {
			return fmt.Errorf("Failed to write line data: %v", err)
		}
	}
	if _, err := w.w.Write([]byte(formatCRC24(w.crc))); nil != err {
		return fmt.Errorf("Failed to write line checksum: %v", err)
	}
	return nil
}

func formatCRC24(hash crc24.Hash24) string {
	return fmt.Sprintf("%06X\n", hash.Sum24())
}

func bufLen(l uint8) int {
	// Reduce by static line overhead: 6 chars for line number and 6 for CRC.
	l -= 12

	// Reduce by the number of whitespaces in between data.
	if l%6 == 0 {
		l -= l/6 - 1
	} else {
		l -= l / 6
	}
	// Reduce by the whitespace between data an CRC.
	l--

	return int(l)
}
