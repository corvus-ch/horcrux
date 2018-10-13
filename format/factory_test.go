package format_test

import (
	"testing"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/logr/buffered"
	"github.com/corvus-ch/horcrux/format/raw"
	"github.com/stretchr/testify/assert"
)

func createTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v.", err)
	}
	return dir
}

func assertTempFile(t *testing.T, path string, size int64) {
	i, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Failed to reade file stats: %v.", err)
	}
	assert.Equal(t, size, i.Size())
}

func TestFactory_Create(t *testing.T) {
	tmp := createTempDir(t)
	defer os.RemoveAll(tmp)
	log := buffered.New(0)

	tests := []struct {
		name      string
		encrypted bool
		assert    func(*testing.T)
	}{
		{"plain",false, func(t *testing.T) {
			assertTempFile(t, filepath.Join(tmp, "plain.raw.042"), 5)
		}},
		{"encrypted", true, func(t *testing.T) {
			assertTempFile(t, filepath.Join(tmp, "encrypted.raw.042"), 100)
			assert.Regexp(t, "INFO Password for 042: [ybndrfg8ejkmcpqxot1uwisza345h769]{12}", log.Buf().String())
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log.Buf().Reset()
			r := raw.New(filepath.Join(tmp, test.name))
			f := format.NewFactory([]format.Format{r}, test.encrypted, log)
			defer func() {
				assert.Nil(t, f.Close())
			}()
			w, err := f.Create(42)
			if err != nil {
				t.Errorf("Failed to create writer: %v.", err)
			}
			w.Write([]byte(test.name))
			test.assert(t)
		})
	}
}
