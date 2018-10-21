package qr_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/corvus-ch/horcrux/format"
	formatAssert "github.com/corvus-ch/horcrux/format/internal/assert"
	"github.com/corvus-ch/horcrux/format/qr"
	"github.com/stretchr/testify/assert"
)

func factory(s string) format.Format {
	return qr.New(s)
}

func TestFormat_Reader(t *testing.T) {
	_, err := factory("").Reader(bytes.NewReader([]byte{}))
	assert.Error(t, err)
}

func TestFormat_Writer(t *testing.T) {
	formatAssert.DataWrite(t, factory, ".png", func(file string, x byte) []string {
		base := filepath.Base(file)
		name := base[0 : len(base)-len(filepath.Ext(base))]
		f, err := os.Stat(file)
		if err != nil {
			return []string{}
		}
		num := (f.Size() / 2120) + 1

		names := make([]string, num)
		for i := int64(0); i < num; i++ {
			names[i] = fmt.Sprintf("%s.%03d.%d.png", name, x, i+1)
		}

		return names
	})
}

func TestFormat_Name(t *testing.T) {
	assert.Equal(t, qr.Name, qr.New("").Name())
}
