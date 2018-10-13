package format

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"crypto/rand"
	"github.com/bketelsen/logr"
	zbase32enc "github.com/corvus-ch/zbase32"
	"golang.org/x/crypto/openpgp"
)

// Factory creates output writers for the different part of a Shamir split part.
type Factory struct {
	io.Closer
	c       []io.Closer // List of Closer instances waiting to be closed.
	formats []Format    // Per format configuration. If nil the format is disabled.
	encrypt bool        // Whether or not the output should be encrypted with a random password
	log     logr.Logger // The output where the auto generated passwords are written to
}

// NewFactory creates a new instance of Factory.
func NewFactory(c []Format, encrypt bool, log logr.Logger) *Factory {
	return &Factory{formats: c, encrypt: encrypt, log: log}
}

// Create is a factory method to be passed to the Shamir splitter.
func (f *Factory) Create(x byte) (io.Writer, error) {
	ws := make([]io.Writer, len(f.formats))

	var i int
	var err error
	for _, format := range f.formats {
		if ws[i], err = f.create(x, format); err != nil {
			return nil, err
		}
		i++
	}

	if !f.encrypt {
		return io.MultiWriter(ws...), nil
	}

	return f.encryptWriter(io.MultiWriter(ws...), x)
}

func (f *Factory) create(x byte, format Format) (io.Writer, error) {
	file, err := os.Create(format.OutputFileName(x))
	if nil != err {
		return nil, err
	}
	f.c = append(f.c, file)
	w, c, err := format.Writer(file)
	if nil != err {
		return nil, err
	}
	if nil != c {
		f.c = append(f.c, c...)
	}

	return w, nil
}

func (f *Factory) encryptWriter(w io.Writer, x byte) (io.Writer, error) {
	p, err := generatePassword()
	if nil != err {
		return nil, err
	}
	hints := &openpgp.FileHints{IsBinary: true, FileName: fmt.Sprintf("%03d.gpg", x)}
	cypher, err := openpgp.SymmetricallyEncrypt(w, p, hints, nil)
	if nil != err {
		return nil, err
	}
	f.c = append(f.c, cypher)

	f.log.Infof("Password for %03d: %s", x, p)

	return cypher, nil
}

// Close closes all the open file handles.
// This is done in reverse order like it would happen by using defer.
func (f *Factory) Close() error {
	for i := len(f.c) - 1; i >= 0; i-- {
		err := f.c[i].Close()
		if nil != err {
			return err
		}
	}
	return nil
}

func generatePassword() ([]byte, error) {
	var buf bytes.Buffer

	enc := zbase32enc.NewEncoder(zbase32enc.StdEncoding, &buf)
	if _, err := io.CopyN(enc, rand.Reader, int64(zbase32enc.StdEncoding.DecodedLen(12))); nil != err {
		return nil, err
	}
	if err := enc.Close(); nil != err {
		return nil, err
	}

	return buf.Bytes(), nil
}
