package create_test

import (
	"bytes"
	"io"

	"github.com/corvus-ch/horcrux/format"
	"github.com/stretchr/testify/mock"
)

type Config struct {
	mock.Mock
}

func NewConfig(input string, f format.Format, encrypt bool) *Config {
	cfg := &Config{}
	cfg.On("Input").Maybe().Return(bytes.NewBufferString(input), nil)
	cfg.On("Formats").Maybe().Return([]format.Format{f}, nil)
	cfg.On("Encrypt").Maybe().Return(encrypt)
	cfg.On("Parts").Maybe().Return(3)
	cfg.On("Threshold").Maybe().Return(2)

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

func (c *Config) Formats() ([]format.Format, error) {
	args := c.Called()
	f := args.Get(0)
	if f == nil {
		return nil, args.Error(1)
	}

	return f.([]format.Format), args.Error(1)
}

func (c *Config) Encrypt() bool {
	args := c.Called()
	return args.Bool(0)
}

func (c *Config) Parts() int {
	args := c.Called()
	return args.Int(0)
}

func (c *Config) Threshold() int {
	args := c.Called()
	return args.Int(0)
}

func (c *Config) LineLength() uint8 {
	args := c.Called()
	return uint8(args.Int(0))
}
