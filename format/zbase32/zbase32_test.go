package zbase32_test

import (
	"github.com/corvus-ch/horcrux/format"
	"testing"

	formatAssert "github.com/corvus-ch/horcrux/format/internal/assert"
	"github.com/corvus-ch/horcrux/format/zbase32"
	"github.com/stretchr/testify/assert"
)

var nameTests = []formatAssert.NameTest{
	{0, "mollis", "mollis.zbase32.000"},
	{1, "commodo", "commodo.zbase32.001"},
	{42, "pellentesque", "pellentesque.zbase32.042"},
	{181, "fringilla", "fringilla.zbase32.181"},
	{254, "venenatis", "venenatis.zbase32.254"},
	{255, "ridiculus", "ridiculus.zbase32.255"},
}

var dataTests = []formatAssert.DataTest{
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

func factory(s string) format.Format {
	return zbase32.New(s)
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
	assert.Equal(t, zbase32.Name, zbase32.New("").Name())
}
