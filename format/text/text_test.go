package text_test

import (
	"testing"

	"github.com/corvus-ch/horcrux/format"
	formatAssert "github.com/corvus-ch/horcrux/format/internal/assert"
	"github.com/corvus-ch/horcrux/format/text"
	"github.com/stretchr/testify/assert"
)

var nameTests = []formatAssert.NameTest{
	{0, "mollis", "mollis.txt.000"},
	{1, "commodo", "commodo.txt.001"},
	{42, "pellentesque", "pellentesque.txt.042"},
	{181, "fringilla", "fringilla.txt.181"},
	{254, "venenatis", "venenatis.txt.254"},
	{255, "ridiculus", "ridiculus.txt.255"},
}

var dataTests = []formatAssert.DataTest{
	{[]byte{240, 191, 199}, "   1: 6n9hq BF444A\n   2: BF444A\n"},
	{[]byte{212, 122, 4}, "   1: 4t7ye 7008A1\n   2: 7008A1\n"},
	{[]byte{0xff}, "   1: 9h 78D648\n   2: 78D648\n"},
	{[]byte{0xb5}, "   1: sw 37CA82\n   2: 37CA82\n"},
	{[]byte{0x34, 0x5a}, "   1: gtpy 358F66\n   2: 358F66\n"},
	{[]byte{0xff, 0xff, 0xff, 0xff, 0xff}, "   1: 99999 999 302370\n   2: 302370\n"},
	{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, "   1: 99999 9999h 9AD7E3\n   2: 9AD7E3\n"},
	{[]byte{
		0xc0, 0x73, 0x62, 0x4a, 0xaf, 0x39, 0x78, 0x51,
		0x4e, 0xf8, 0x44, 0x3b, 0xb2, 0xa8, 0x59, 0xc7,
		0x5f, 0xc3, 0xcc, 0x6a, 0xf2, 0x6d, 0x5a, 0xaa,
	}, "   1: ab3sr 1ix8f hfnuz aeo75 fkn3a 7xh8u dk6js iiko 039DEA\n   2: 039DEA\n"},
}

func factory(s string) format.Format {
	return text.New(s)
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
	assert.Equal(t, text.Name, text.New("").Name())
}
