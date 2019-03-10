package internal

import (
	"bytes"
	"io"
	"os"

	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/input"
	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Config struct {
	*viper.Viper

	input input.Input
}

func NewConfig() (*Config, error) {
	v := viper.New()

	setDefaults(v)

	v.SetConfigName(".horcrux")
	v.AddConfigPath("$HOME/")
	v.AddConfigPath(".")

	err := v.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); err != nil && !ok {
		return nil, err
	}

	return &Config{Viper: v}, nil
}

func NewConfigFromYaml(yaml string) (*Config, error) {
	v := viper.New()

	setDefaults(v)

	v.SetConfigType("yaml")

	return &Config{Viper: v}, v.ReadConfig(bytes.NewBuffer([]byte(yaml)))
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("parts", 3)
	v.SetDefault("threshold", 2)
	v.SetDefault("format", "text")
}

func (c *Config) Encrypted() bool {
	return c.GetBool("encrypted")
}

func (c *Config) Decrypt() bool {
	return c.Encrypted()
}

func (c *Config) Format() (format.Format, error) {
	return format.New(c.GetString("format"), c.InputInfo())
}

func (c *Config) Formats() ([]format.Format, error) {
	formatNames := c.GetStringSlice("formats")
	formats := make([]format.Format, len(formatNames))

	for i, f := range formatNames {
		ff, err := format.New(f, c.InputInfo())
		if err != nil {
			return nil, err
		}
		formats[i] = ff
	}

	return formats, nil
}

func (c *Config) Reader() (io.Reader, error) {
	name := c.GetString("input")
	if name == "-" || name == "" {
		return os.Stdin, nil
	}

	return os.Open(name)
}

func (c *Config) InputInfo() input.Input {
	if c.input != nil {
		return c.input
	}

	name := c.GetString("input")
	if name == "-" || name == "" {
		c.input = input.NewStreamInput(c.GetString("stem"))
	} else {

		c.input = input.NewFileInput(name, c.GetString("stem"))
	}

	return c.input
}

func (c *Config) Output() (io.Writer, error) {
	out := c.GetString("output")
	if len(out) == 0 || out == "-" {
		return os.Stdout, nil
	}

	return os.Create(out)
}

func (c *Config) Parts() int {
	return c.GetInt("parts")
}

func (c *Config) Threshold() int {
	return c.GetInt("threshold")
}

func (c *Config) FileNames() []string {
	return c.GetStringSlice("files")
}

func (c *Config) StringValue(key string) kingpin.Value {
	return &stringValue{value{key: key, viper: c.Viper}}
}

func (c *Config) StringsValue(key string) kingpin.Value {
	return &stringsValue{value{key: key, viper: c.Viper}}
}

func (c *Config) IntValue(key string) kingpin.Value {
	return &intValue{value{key: key, viper: c.Viper}}
}

func (c *Config) BoolValue(key string) kingpin.Value {
	return &boolValue{value{key: key, viper: c.Viper}}
}

func (c *Config) EnumValue(key string, options ...string) kingpin.Value {
	return &enumValue{value: value{key: key, viper: c.Viper}, options: options}
}

func (c *Config) EnumsValue(key string, options ...string) kingpin.Value {
	return &enumsValue{value: value{key: key, viper: c.Viper}, options: options}
}
