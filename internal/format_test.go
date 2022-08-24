package internal_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bketelsen/logr"
	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/format/raw"
	"github.com/corvus-ch/horcrux/format/text"
	"github.com/corvus-ch/horcrux/format/zbase32"
	"github.com/corvus-ch/horcrux/internal"
	"github.com/corvus-ch/logr/buffered"
	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
)

func assertFormatAction(t *testing.T, args []string, action func(format.Config, logr.Logger) error) {
	log := buffered.New(0)
	app := kingpin.New("test", "test")
	internal.RegisterFormatCommand(app, log, action)
	_, err := app.Parse(args)
	assert.Nil(t, err)
}

func TestFormatCommand_Input(t *testing.T) {
	file, err := ioutil.TempFile("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	tests := []struct {
		name string
		args []string
		file *os.File
	}{
		{"default", []string{"format"}, os.Stdin},
		{"stdin", []string{"format", "--", "-"}, os.Stdin},
		{"file", []string{"format", "--", file.Name()}, file},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertFormatAction(t, test.args, func(cfg format.Config, _ logr.Logger) error {
				reader, err := cfg.Input()
				assert.Nil(t, err)
				assert.Equal(t, test.file.Name(), reader.(*os.File).Name())
				return nil
			})
		})
	}
}

func TestFormatCommand_Formats(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		formats []string
	}{
		{"default", []string{"format"}, []string{text.Name}},
		{"single", []string{"format", "-f", "raw"}, []string{raw.Name}},
		{"multiple", []string{"format", "-f", "raw", "-f", "zbase32"}, []string{raw.Name, zbase32.Name}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertFormatAction(t, test.args, func(cfg format.Config, _ logr.Logger) error {
				formats, err := cfg.Formats()
				if err != nil {
					t.Fatal(err)
				}
				if len(test.formats) != len(formats) {
					t.Fatalf("expected %d formats, got %d", len(test.formats), len(formats))
				}
				for i, name := range test.formats {
					assert.Equal(t, name, formats[i].Name())
				}
				return nil
			})
		})
	}
}
