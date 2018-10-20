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

func factory(s string) format.Format {
	return text.New(s)
}

func TestFormat_OutputFileName(t *testing.T) {
	formatAssert.Name(t, nameTests, factory)
}

func TestFormat_Reader(t *testing.T) {
	formatAssert.DataRead(t, factory, ".txt")
}

func TestFormat_Writer(t *testing.T) {
	formatAssert.DataWrite(t, factory, ".txt")
}

func TestFormat_Name(t *testing.T) {
	assert.Equal(t, text.Name, text.New("").Name())
}
