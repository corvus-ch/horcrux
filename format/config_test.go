package format_test

import (
	"bytes"
	"io"

	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/input"
	"github.com/stretchr/testify/mock"
)

type Config struct {
	mock.Mock
}

func NewConfig(name string, f format.Format) *Config {
	cfg := &Config{}
	cfg.On("Input").Maybe().Return(bytes.NewBufferString(name), nil)
	cfg.On("InputInfo").Maybe().Return(input.NewStreamInput(""))
	cfg.On("Formats").Maybe().Return([]format.Format{f}, nil)

	return cfg
}

func (c *Config) Input() (io.Reader, error) {
	args := c.Called()
	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}

	return r.(io.Reader), args.Error(1)
}

func (c *Config) InputInfo() input.Input {
	args := c.Called()

	return args.Get(0).(input.Input)
}

func (c *Config) Formats() ([]format.Format, error) {
	args := c.Called()
	f := args.Get(0)
	if f == nil {
		return nil, args.Error(1)
	}

	return f.([]format.Format), args.Error(1)
}
