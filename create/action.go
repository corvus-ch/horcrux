package create

import (
	"fmt"
	"io"

	"github.com/bketelsen/logr"
	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/shamir"
)

// Create creates a new set of secret parts according to the given config.
func Create(cfg Config, log logr.Logger) (result error) {
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
	defer func() {
		if err := factory.Close(); err != nil && result == nil {
			result = err
		}
	}()

	writer, err := shamir.NewWriter(cfg.Parts(), cfg.Threshold(), factory.Create)
	if nil != err {
		return fmt.Errorf("failed to create processing pipeline: %v", err)
	}

	if _, err := io.Copy(writer, reader); nil != err {
		result = fmt.Errorf("failed to process data: %v", err)
	}

	return
}
