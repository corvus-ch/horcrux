package text

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"text/template"

	"github.com/corvus-ch/horcrux/input"
	"github.com/martinlindhe/crc24"
	"github.com/masterminds/sprig"
)

const defaultTemplate = `
{{- block "header" .}}{{end -}}
{{- block "lines" .Lines -}}
{{range .}}
  {{- block "line" . -}}
    {{printf "% 4d" .Number }}: {{if .Data}}{{ .Data }} {{end}}{{printf "%06X" .CRC }}
  {{- end}}
{{end}}
{{- end -}}
{{- block "footer" .}}{{end -}}
`

type writer struct {
	io.WriteCloser
	err  error
	w    io.Writer    // The underlying writer
	n    int          // Lines count
	buf  []byte       // buffered data waiting to be encoded
	nbuf int          // number of bytes in buf
	crc  crc24.Hash24 // The checksum
	t    *template.Template
	data templateData
	wg   sync.WaitGroup
}

type line struct {
	Number int
	Data   string
	CRC    uint32
}

type output struct {
	Name string
	X    byte
}

type templateData struct {
	Input  input.Input
	Lines  chan line
	Output output
}

// NewWriter returns an text format writer instance.
func NewWriter(w io.Writer, f *Format, o output) (io.WriteCloser, error) {
	tw := &writer{
		w:   w,
		buf: make([]byte, bufLen(f.LineLength)),
		crc: crc24.New(),
		data: templateData{
			Input:  f.input,
			Lines:  make(chan line),
			Output: o,
		},
		t: template.New("text"),
	}

	tw.t.Funcs(sprig.TxtFuncMap())

	if err := tw.parse(); err != nil {
		return nil, err
	}

	go tw.render()

	return tw, nil
}

func (w *writer) parse() error {
	if _, err := w.t.Parse(defaultTemplate); err != nil {
		return err
	}

	if _, err := w.t.ParseGlob("text.tmpl"); err != nil && !strings.Contains(err.Error(), "pattern matches no files") {
		return err
	}

	return nil
}

func (w *writer) render() {
	w.wg.Add(1)
	w.err = w.t.Execute(w.w, w.data)
	w.wg.Done()
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
		w.writeLine(w.buf)
		w.nbuf = 0
	}

	// Large interior chunks.
	for len(p) > len(w.buf) {
		nn := len(w.buf)
		w.crc.Write(p[:nn])
		w.writeLine(p[:nn])
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
		w.writeLine(w.buf[:w.nbuf])
	}
	w.writeLine(nil)
	close(w.data.Lines)
	w.wg.Wait()
	return w.err
}

func (w *writer) writeLine(data []byte) {
	w.n++
	var buf bytes.Buffer
	for i, r := range data {
		buf.WriteByte(r)
		if i%5 == 4 && i != len(data)-1 {
			buf.WriteRune(' ')
		}
	}
	w.data.Lines <- line{Number: w.n, Data: buf.String(), CRC: w.crc.Sum24()}
}

func bufLen(l uint8) int {
	// Reduce by static Line overhead: 6 chars for Line number and 6 for CRC.
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
