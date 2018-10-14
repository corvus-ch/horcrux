package restore_test

import (
	"bytes"
	"io"

	"github.com/corvus-ch/horcrux/format"
	"github.com/stretchr/testify/mock"
)

type Config struct {
	mock.Mock
}

func NewConfig(files []string, format format.Format, decrypt bool) (*Config, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	cfg := &Config{}
	cfg.On("Format").Maybe().Return(format, nil)
	cfg.On("Output").Maybe().Return(buf, nil)
	cfg.On("FileNames").Maybe().Return(files)
	cfg.On("Decrypt").Maybe().Return(decrypt)

	return cfg, buf
}

func (c *Config) Format() (format.Format, error) {
	args := c.Called()
	f := args.Get(0)
	if f == nil {
		return nil, args.Error(1)
	}

	return f.(format.Format), args.Error(1)
}

func (c *Config) Output() (io.Writer, error) {
	args := c.Called()
	w := args.Get(0)
	if w == nil {
		return nil, args.Error(1)
	}

	return w.(io.Writer), args.Error(1)
}

func (c *Config) FileNames() []string {
	args := c.Called()
	return args.Get(0).([]string)
}

func (c *Config) Decrypt() bool {
	args := c.Called()
	return args.Bool(0)
}
