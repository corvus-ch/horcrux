package assert

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/input"
	"github.com/sebdah/goldie"
	"github.com/stretchr/testify/assert"
)

const x = byte(42)

// FormatFactory describes a func used for instantiating a Format during assertions.
type FormatFactory func(input.Input) format.Format

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
			r, err := factory(input.NewStreamInput("")).Reader(f)
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
			f, err := os.Open(file)
			if err != nil {
				t.Fatal(err)
			}
			dir, err := ioutil.TempDir("", name)
			defer os.RemoveAll(dir)
			if err != nil {
				t.Fatal(err)
			}

			w, cl, err := factory(newInputMock(t, f, filepath.Join(dir, name))).Writer(x)
			assert.NoError(t, err)
			io.Copy(w, f)
			assert.NoError(t, err)
			assertClose(t, cl)
			for _, outfile := range outfilenames(file, x) {
				assertFileContent(t, dir, outfile, suffix)
			}
		})
	}
}

func newInputMock(t *testing.T, file *os.File, stem string) input.Input {
	return input.NewFileInput(file.Name(), stem)
}

func assertFileContent(t *testing.T, dir, name, suffix string) {
	file, err := os.Open(filepath.Join(dir, name))
	if err != nil {
		t.Error(err)
		return
	}
	out, _ := ioutil.ReadAll(file)
	goldenName := strings.Replace(strings.Replace(name, suffix, "", -1), fmt.Sprintf(".%03d", x), "", -1)
	goldie.Assert(t, goldenName, out)
}

func assertClose(t *testing.T, cs []io.Closer) {
	for i := len(cs); i > 0; i-- {
		assert.NoError(t, cs[i-1].Close())
	}
}
