package raw_test

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"testing"

	"github.com/corvus-ch/horcrux/format"
	formatAssert "github.com/corvus-ch/horcrux/format/internal/assert"
	"github.com/corvus-ch/horcrux/format/raw"
	"github.com/stretchr/testify/assert"
)

var nameTests = []formatAssert.NameTest{
	{0, "mollis", "mollis.raw.000"},
	{1, "commodo", "commodo.raw.001"},
	{42, "pellentesque", "pellentesque.raw.042"},
	{181, "fringilla", "fringilla.raw.181"},
	{254, "venenatis", "venenatis.raw.254"},
	{255, "ridiculus", "ridiculus.raw.255"},
}

func factory(s string) format.Format {
	return raw.New(s)
}

func TestFormat_OutputFileName(t *testing.T) {
	formatAssert.Name(t, nameTests, factory)
}

func TestFormat_Reader(t *testing.T) {
	data := randomData()
	r, err := raw.New("").Reader(bytes.NewBuffer(data))
	assert.Nil(t, err)
	out, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, data, out)
}

func TestFormat_Writer(t *testing.T) {
	out := &bytes.Buffer{}
	w, cl, err := raw.New("").Writer(out)
	assert.Nil(t, cl)
	assert.NoError(t, err)
	data := randomData()
	w.Write(data)
	assert.Equal(t, data, out.Bytes())
}

func TestFormat_Name(t *testing.T) {
	assert.Equal(t, raw.Name, raw.New("").Name())
}

func randomData() []byte {
	data := make([]byte, rand.Intn(1024))
	rand.Read(data)

	return data
}
