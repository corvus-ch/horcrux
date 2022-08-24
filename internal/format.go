package internal

import (
	"io"
	"os"

	"github.com/bketelsen/logr"
	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/input"
	"gopkg.in/alecthomas/kingpin.v2"
)

type formatAction func(cfg format.Config, logger logr.Logger) error

type formatCommand struct {
	action formatAction

	// The logger
	log logr.Logger

	// Arguments
	input string

	// Flags
	formats  []string
	stemFlag string

	// internal
	info input.Input
}

// RegisterFormatCommand registers the format sub command with the application.
func RegisterFormatCommand(app *kingpin.Application, log logr.Logger, action formatAction) *formatCommand {
	c := &formatCommand{action: action, log: log}

	cc := app.Command("format", "formats input the same way as create would do without doing the split")
	cc.Action(c.Execute)
	cc.Arg("input", "the input file").
		StringVar(&c.input)
	cc.Flag("format", "the output formats").
		Default(format.Default).
		Short('f').
		StringsVar(&c.formats)
	cc.Flag("output", "name stem for the output files").
		Short('o').
		StringVar(&c.stemFlag)

	return c
}

// Execute runs the action callback.
func (c *formatCommand) Execute(_ *kingpin.ParseContext) error {
	return c.action(c, c.log)
}

// Input returns the input file.
func (c *formatCommand) Input() (io.Reader, error) {
	if c.input == "-" || c.input == "" {
		return os.Stdin, nil
	}

	return os.Open(c.input)
}

// InputInfo returns detail infor about the input.
func (c *formatCommand) InputInfo() input.Input {
	if c.info != nil {
		return c.info
	}

	if c.input == "-" || c.input == "" {
		c.info = input.NewStreamInput(c.stemFlag)
	} else {

		c.info = input.NewFileInput(c.input, c.stemFlag)
	}

	return c.info
}

// Formats returns the list of output formats to produce
func (c *formatCommand) Formats() ([]format.Format, error) {
	formats := make([]format.Format, len(c.formats))

	for i, f := range c.formats {
		ff, err := format.New(f, c.InputInfo())
		if err != nil {
			return nil, err
		}
		formats[i] = ff
	}

	return formats, nil
}
