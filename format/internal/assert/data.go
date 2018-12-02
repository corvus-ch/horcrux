package assert

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/input"
	"github.com/corvus-ch/horcrux/output"
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

			w, cl, err := factory(newInputMock(t, f, filepath.Join(dir, name))).Writer(x, output.NewOutput())
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
	actualFile, err := os.Open(filepath.Join(dir, name))
	if err != nil {
		t.Error(err)
		return
	}
	goldieName := strings.Replace(strings.Replace(name, suffix, "", -1), fmt.Sprintf(".%03d", x), "", -1)

	if updateFlag := flag.Lookup(goldie.FlagName); updateFlag != nil && updateFlag.Value.String() == "true" {
		actualData, _ := ioutil.ReadAll(actualFile)
		err := goldie.Update(goldieName, actualData)
		if err != nil {
			t.Error(err)
		}
		return
	}

	expectedFile, err := os.Open(filepath.Join(goldie.FixtureDir, fmt.Sprintf("%s%s", goldieName, goldie.FileNameSuffix)))
	if err != nil {
		t.Error(err)
		return
	}

	for {
		if err := compareFileChunk(actualFile, expectedFile); err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}
	}
}

func compareFileChunk(actualFile, expectedFile *os.File) error {
	actualBytes := make([]byte, 512)
	_, actualErr := actualFile.Read(actualBytes)

	expectedBytes := make([]byte, 512)
	_, expectedErr := expectedFile.Read(expectedBytes)

	if actualErr == io.EOF && expectedErr == io.EOF {
		return io.EOF
	} else if actualErr == io.EOF || expectedErr == io.EOF {
		return fmt.Errorf("files have unequal length")
	} else if actualErr != nil || expectedErr != nil {
		return fmt.Errorf("error while reading bytes: %v or %v", actualFile, expectedErr)
	}

	if !bytes.Equal(actualBytes, expectedBytes) {
		return fmt.Errorf("files are not equal")
	}

	return nil
}

func assertClose(t *testing.T, cs []io.Closer) {
	for i := len(cs); i > 0; i-- {
		assert.NoError(t, cs[i-1].Close())
	}
}
