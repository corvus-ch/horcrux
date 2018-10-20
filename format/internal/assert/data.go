package assert

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/corvus-ch/horcrux/format"
	"github.com/sebdah/goldie"
	"github.com/stretchr/testify/assert"
)

// FormatFactory describes a func used for instantiating a Format during assertions.
type FormatFactory func(string) format.Format

// DataRead asserts a formats read behaviour.
func DataRead(t *testing.T, factory FormatFactory, suffix string) {
	goldie.FileNameSuffix = ".bin"
	files, err := filepath.Glob(filepath.Join(goldie.FixtureDir, "*"+suffix))
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		basename := filepath.Base(file)
		name := strings.TrimSuffix(basename, filepath.Ext(basename))
		t.Run(name, func(t *testing.T) {
			f, err := os.Open(file)
			if err != nil {
				t.Fatal(err)
			}
			r, err := factory("").Reader(f)
			assert.Nil(t, err)
			out, err := ioutil.ReadAll(r)
			assert.NoError(t, err)
			goldie.Assert(t, name, out)
		})
	}
}

// DataWrite asserts a formats write behaviour.
func DataWrite(t *testing.T, factory FormatFactory, suffix string) {
	goldie.FileNameSuffix = suffix
	files, err := filepath.Glob(filepath.Join(goldie.FixtureDir, "*.bin"))
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		basename := filepath.Base(file)
		name := strings.TrimSuffix(basename, filepath.Ext(basename))
		t.Run(name, func(t *testing.T) {
			out := &bytes.Buffer{}
			f, err := os.Open(file)
			if err != nil {
				t.Fatal(err)
			}
			w, cl, err := factory("").Writer(out)
			assert.NoError(t, err)
			_, err = io.Copy(w, f)
			assert.NoError(t, err)
			for i := len(cl); i > 0; i-- {
				err = cl[i-1].Close()
				assert.NoError(t, err)
			}
			goldie.Assert(t, name, out.Bytes())
		})
	}
}
