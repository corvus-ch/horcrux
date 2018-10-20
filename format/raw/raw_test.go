package raw_test

import (
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
	formatAssert.DataRead(t, factory, ".raw")
}

func TestFormat_Writer(t *testing.T) {
	formatAssert.DataWrite(t, factory, ".raw")
}

func TestFormat_Name(t *testing.T) {
	assert.Equal(t, raw.Name, raw.New("").Name())
}
