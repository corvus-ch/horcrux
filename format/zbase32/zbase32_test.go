package zbase32_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/corvus-ch/horcrux/format/zbase32"
	"github.com/stretchr/testify/assert"
)

var nameTests = []struct {
	x        byte
	stem     string
	expected string
}{
	{0, "mollis", "mollis.zbase32.000"},
	{1, "commodo", "commodo.zbase32.001"},
	{42, "pellentesque", "pellentesque.zbase32.042"},
	{181, "fringilla", "fringilla.zbase32.181"},
	{254, "venenatis", "venenatis.zbase32.254"},
	{255, "ridiculus", "ridiculus.zbase32.255"},
}

func TestFormat_OutputFileName(t *testing.T) {
	for _, test := range nameTests {
		t.Run(fmt.Sprint(test.x), func(t *testing.T) {
			f := zbase32.New(test.stem)
			assert.Equal(t, test.expected, f.OutputFileName(test.x))
		})
	}
}

var testData = []struct {
	decoded []byte
	encoded string
}{
	{[]byte{240, 191, 199}, "6n9hq"},
	{[]byte{212, 122, 4}, "4t7ye"},
	{[]byte{0xff}, "9h"},
	{[]byte{0xb5}, "sw"},
	{[]byte{0x34, 0x5a}, "gtpy"},
	{[]byte{0xff, 0xff, 0xff, 0xff, 0xff}, "99999999"},
	{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, "999999999h"},
	{[]byte{
		0xc0, 0x73, 0x62, 0x4a, 0xaf, 0x39, 0x78, 0x51,
		0x4e, 0xf8, 0x44, 0x3b, 0xb2, 0xa8, 0x59, 0xc7,
		0x5f, 0xc3, 0xcc, 0x6a, 0xf2, 0x6d, 0x5a, 0xaa,
	}, "ab3sr1ix8fhfnuzaeo75fkn3a7xh8udk6jsiiko"},
}

func TestFormat_Reader(t *testing.T) {
	for _, test := range testData {
		t.Run(test.encoded, func(t *testing.T) {
			r, err := zbase32.New("").Reader(bytes.NewBufferString(test.encoded))
			assert.Nil(t, err)
			out, err := ioutil.ReadAll(r)
			assert.NoError(t, err)
			assert.Equal(t, test.decoded, out)
		})
	}
}

func TestFormat_Writer(t *testing.T) {
	for _, test := range testData {
		t.Run(test.encoded, func(t *testing.T) {
			out := &bytes.Buffer{}
			w, cl, err := zbase32.New("").Writer(out)
			assert.Len(t, cl, 1)
			assert.NoError(t, err)
			w.Write(test.decoded)
			cl[0].Close()
			assert.Equal(t, test.encoded, out.String())
		})
	}
}

func TestFormat_Name(t *testing.T) {
	assert.Equal(t, zbase32.Name, zbase32.New("").Name())
}
