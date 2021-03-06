package internal_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bketelsen/logr"
	"github.com/corvus-ch/horcrux/create"
	"github.com/corvus-ch/horcrux/format/raw"
	"github.com/corvus-ch/horcrux/format/text"
	"github.com/corvus-ch/horcrux/format/zbase32"
	"github.com/corvus-ch/horcrux/internal"
	"github.com/corvus-ch/logr/buffered"
	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
)

func assertCreateAction(t *testing.T, args []string, action func(create.Config, logr.Logger) error) {
	log := buffered.New(0)
	app := kingpin.New("test", "test")
	internal.RegisterCreateCommand(app, log, action)
	_, err := app.Parse(args)
	assert.Nil(t, err)
}

func TestCreateCommand_Encrypt(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		encrypted bool
	}{
		{"plain", []string{"create"}, false},
		{"encrypted", []string{"create", "-e"}, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertCreateAction(t, test.args, func(cfg create.Config, _ logr.Logger) error {
				assert.Equal(t, test.encrypted, cfg.Encrypt())
				return nil
			})
		})
	}
}

func TestCreateCommand_Threshold(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		threshold int
	}{
		{"default", []string{"create"}, 2},
		{"flagged", []string{"create", "-n", "42"}, 42},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertCreateAction(t, test.args, func(cfg create.Config, _ logr.Logger) error {
				assert.Equal(t, test.threshold, cfg.Threshold())
				return nil
			})
		})
	}
}

func TestCreateCommand_Parts(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		parts int
	}{
		{"default", []string{"create"}, 3},
		{"flagged", []string{"create", "-m", "42"}, 42},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertCreateAction(t, test.args, func(cfg create.Config, _ logr.Logger) error {
				assert.Equal(t, test.parts, cfg.Parts())
				return nil
			})
		})
	}
}

func TestCreateCommand_Input(t *testing.T) {
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
		{"default", []string{"create"}, os.Stdin},
		{"stdin", []string{"create", "--", "-"}, os.Stdin},
		{"file", []string{"create", "--", file.Name()}, file},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertCreateAction(t, test.args, func(cfg create.Config, _ logr.Logger) error {
				reader, err := cfg.Input()
				assert.Nil(t, err)
				assert.Equal(t, test.file.Name(), reader.(*os.File).Name())
				return nil
			})
		})
	}
}

func TestCreateCommand_Formats(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		formats []string
	}{
		{"default", []string{"create"}, []string{text.Name}},
		{"single", []string{"create", "-f", "raw"}, []string{raw.Name}},
		{"multiple", []string{"create", "-f", "raw", "-f", "zbase32"}, []string{raw.Name, zbase32.Name}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertCreateAction(t, test.args, func(cfg create.Config, _ logr.Logger) error {
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
