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

func factory(s string) format.Format {
	return zbase32.New(s)
}

func TestFormat_OutputFileName(t *testing.T) {
	formatAssert.Name(t, nameTests, factory)
}

func TestFormat_Reader(t *testing.T) {
	formatAssert.DataRead(t, factory, ".zbase32")
}

func TestFormat_Writer(t *testing.T) {
	formatAssert.DataWrite(t, factory, ".zbase32")
}

func TestFormat_Name(t *testing.T) {
	assert.Equal(t, zbase32.Name, zbase32.New("").Name())
}
