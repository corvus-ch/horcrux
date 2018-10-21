package assert

import (
	"crypto/rand"
	"fmt"
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

// OutputFileNames describes a func used to get the output file names only known to the calling test case.
type OutputFileNames func(file string, x byte) []string

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
func DataWrite(t *testing.T, factory FormatFactory, suffix string, outfilenames OutputFileNames) {
	goldie.FileNameSuffix = suffix
	files, err := filepath.Glob(filepath.Join(goldie.FixtureDir, "*.bin"))
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		basename := filepath.Base(file)
		name := strings.TrimSuffix(basename, filepath.Ext(basename))
		t.Run(name, func(t *testing.T) {
			x := randomByte(t)
			f, err := os.Open(file)
			if err != nil {
				t.Fatal(err)
			}
			dir, err := ioutil.TempDir("", name)
			defer os.RemoveAll(dir)
			if err != nil {
				t.Fatal(err)
			}
			subject := factory(filepath.Join(dir, name))
			w, cl, err := subject.Writer(x)
			assert.NoError(t, err)
			_, err = io.Copy(w, f)
			assert.NoError(t, err)
			for i := len(cl); i > 0; i-- {
				err = cl[i-1].Close()
				assert.NoError(t, err)
			}
			for _, outfile := range outfilenames(file, x) {
				file, err := os.Open(filepath.Join(dir, outfile))
				if err != nil {
					t.Fatal(err)
				}
				out, _ := ioutil.ReadAll(file)
				goldenName := strings.Replace(strings.Replace(outfile, suffix, "", -1), fmt.Sprintf(".%03d", x), "", -1)
				t.Log(goldenName)
				goldie.Assert(t, goldenName, out)
			}
		})
	}
}

func randomByte(t *testing.T) byte {
	b := make([]byte, 1)
	if _, err := rand.Read(b); err != nil {
		t.Fatal(err)
	}

	return b[0]
}
