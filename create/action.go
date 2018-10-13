package create

import (
	"fmt"
	"io"

	"github.com/bketelsen/logr"
	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/shamir"
)

func Create(cfg Config, log logr.Logger) error {
	reader, err := cfg.Input()
	if err != nil {
		return fmt.Errorf("failed to open input: %v", err)
	}

	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	formats, err := cfg.Formats()
	if err != nil {
		return fmt.Errorf("failed to setup output formatting: %v", err)
	}

	factory := format.NewFactory(formats, cfg.Encrypt(), log)

	writer, err := shamir.NewWriter(cfg.Parts(), cfg.Threshold(), factory.Create)
	if nil != err {
		return fmt.Errorf("failed to create processing pipeline: %v", err)
	}

	if _, err := io.Copy(writer, reader); nil != err {
		return fmt.Errorf("failed to process data: %v", err)
	}

	if err := factory.Close(); nil != err {
		return fmt.Errorf("failed to close open files: %v", err)
	}

	return nil
}
