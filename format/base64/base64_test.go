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

func factory(s string) format.Format {
	return base64.New(s)
}

func TestFormat_OutputFileName(t *testing.T) {
	formatAssert.Name(t, nameTests, factory)
}

func TestFormat_Reader(t *testing.T) {
	formatAssert.DataRead(t, factory, ".base64")
}

func TestFormat_Writer(t *testing.T) {
	formatAssert.DataWrite(t, factory, ".base64")
}

func TestFormat_Name(t *testing.T) {
	assert.Equal(t, base64.Name, base64.New("").Name())
}
