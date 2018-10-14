package text

import (
	"bytes"
	"fmt"
	"io"

	"github.com/martinlindhe/crc24"
)

type reader struct {
	io.Reader
	p   *parser      // The underlying parser
	crc crc24.Hash24 // The checksum of the whole document.
	buf bytes.Buffer // buffered data waiting to read.
	eof bool
}

func NewReader(r io.Reader) io.Reader {
	return &reader{p: NewParser(r), crc: crc24.New()}
}

func (r *reader) Read(p []byte) (int, error) {
	var n int
	if r.eof {
		return 0, io.EOF
	}

	if r.buf.Len() == 0 && !r.eof {
		data, err := r.readLine()
		if nil != err && io.EOF != err {
			return n, err
		} else if io.EOF == err {
			r.eof = true
		}

		r.buf.Write(data)
	}

	for n < len(p) && r.buf.Len() > 0 {
		n += copy(p[n:], r.buf.Next(min(r.buf.Len(), len(p))))
	}

	if r.eof == true && r.buf.Len() == 0 {
		return n, io.EOF
	}

	return n, nil
}

func (r *reader) readLine() ([]byte, error) {
	l, err := r.p.Parse()
	if nil != err && io.EOF != err {
		return nil, err
	}

	if nil == l {
		return nil, io.EOF
	}

	data := []byte(l.Data)

	if _, err := r.crc.Write(data); nil != err {
		return nil, err
	}

	if r.crc.Sum24() != l.CRC {
		return nil, fmt.Errorf(`checksum does not match on line %v: %06X`, l.N, r.crc.Sum24())
	}

	// Empty line means we have reached the last line containing the document checksum.
	if len(data) == 0 {
		return nil, io.EOF
	}

	return data, nil
}
