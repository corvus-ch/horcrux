package qr_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/corvus-ch/horcrux/format"
	formatAssert "github.com/corvus-ch/horcrux/format/internal/assert"
	"github.com/corvus-ch/horcrux/format/qr"
	"github.com/stretchr/testify/assert"
)

func factory(s string) format.Format {
	return qr.New(s)
}

func TestFormat_Reader(t *testing.T) {
	_, err := factory("").Reader(bytes.NewReader([]byte{}))
	assert.Error(t, err)
}

func TestFormat_Writer(t *testing.T) {
	formatAssert.DataWrite(t, factory, ".png")

	t.Run("too much data", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "fail")
		defer os.RemoveAll(dir)
		w, cl, err := factory(filepath.Join(dir, "fail")).Writer(255)
		assert.NoError(t, err)
		r := io.LimitReader(rand.Reader, 2120)
		_, err = io.Copy(w, r)
		for i := len(cl); i > 0; i-- {
			e := cl[i-1].Close()
			if e != nil {
				err = e
			}
		}
		assert.Error(t, err)
	})
}

func TestFormat_Name(t *testing.T) {
	assert.Equal(t, qr.Name, qr.New("").Name())
}
