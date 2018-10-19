package base64_test

import (
	"testing"

	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/format/base64"
	formatAssert "github.com/corvus-ch/horcrux/format/internal/assert"
	"github.com/stretchr/testify/assert"
)

var nameTests = []formatAssert.NameTest{
	{0, "mollis", "mollis.base64.000"},
	{1, "commodo", "commodo.base64.001"},
	{42, "pellentesque", "pellentesque.base64.042"},
	{181, "fringilla", "fringilla.base64.181"},
	{254, "venenatis", "venenatis.base64.254"},
	{255, "ridiculus", "ridiculus.base64.255"},
}

var dataTests = []formatAssert.DataTest{
	{[]byte{0}, "AA=="},
	{[]byte{0xff}, "/w=="},
	{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, "////////"},
	{[]byte{
		0xc0, 0x73, 0x62, 0x4a, 0xaf, 0x39, 0x78, 0x51,
		0x4e, 0xf8, 0x44, 0x3b, 0xb2, 0xa8, 0x59, 0xc7,
		0x5f, 0xc3, 0xcc, 0x6a, 0xf2, 0x6d, 0x5a, 0xaa,
	}, "wHNiSq85eFFO+EQ7sqhZx1/DzGrybVqq"},
}

func factory(s string) format.Format {
	return base64.New(s)
}

func TestFormat_OutputFileName(t *testing.T) {
	formatAssert.Name(t, nameTests, factory)
}

func TestFormat_Reader(t *testing.T) {
	formatAssert.DataRead(t, dataTests, factory)
}

func TestFormat_Writer(t *testing.T) {
	formatAssert.DataWrite(t, dataTests, factory)
}

func TestFormat_Name(t *testing.T) {
	assert.Equal(t, base64.Name, base64.New("").Name())
}
