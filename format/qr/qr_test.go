package qr_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"io/ioutil"
	"testing"

	"github.com/corvus-ch/horcrux/format"
	formatAssert "github.com/corvus-ch/horcrux/format/internal/assert"
	"github.com/corvus-ch/horcrux/format/qr"
	"github.com/stretchr/testify/assert"
)

var nameTests = []formatAssert.NameTest{
	{0, "mollis", "mollis.000.png"},
	{1, "commodo", "commodo.001.png"},
	{42, "pellentesque", "pellentesque.042.png"},
	{181, "fringilla", "fringilla.181.png"},
	{254, "venenatis", "venenatis.254.png"},
	{255, "ridiculus", "ridiculus.255.png"},
}

func factory(s string) format.Format {
	return qr.New(s)
}

func TestFormat_OutputFileName(t *testing.T) {
	formatAssert.Name(t, nameTests, factory)
}

func TestFormat_Reader(t *testing.T) {
	_, err := factory("").Reader(bytes.NewReader([]byte{}))
	assert.Error(t, err)
}

func TestFormat_Writer(t *testing.T) {
	formatAssert.DataWrite(t, factory, ".png")

	t.Run("too much data", func(t *testing.T) {
		w, cl, err := factory("").Writer(ioutil.Discard)
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
