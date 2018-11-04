package raw_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/corvus-ch/horcrux/format"
	formatAssert "github.com/corvus-ch/horcrux/format/internal/assert"
	"github.com/corvus-ch/horcrux/format/raw"
	"github.com/corvus-ch/horcrux/input"
	"github.com/stretchr/testify/assert"
)

func factory(i input.Input) format.Format {
	return raw.New(i)
}

func TestFormat_Reader(t *testing.T) {
	formatAssert.DataRead(t, factory, ".raw")
}

func TestFormat_Writer(t *testing.T) {
	formatAssert.DataWrite(t, factory, ".raw", func(file string, x byte) []string {
		base := filepath.Base(file)
		name := base[0 : len(base)-len(filepath.Ext(base))]
		return []string{fmt.Sprintf("%s.raw.%03d", name, x)}
	})
}

func TestFormat_Name(t *testing.T) {
	assert.Equal(t, raw.Name, raw.New(nil).Name())
}
