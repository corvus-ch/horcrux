package format

import (
	"fmt"
	"io"

	"github.com/bketelsen/logr"
)

// DoFormat applies a format to a given input.
func DoFormat(cfg Config, log logr.Logger) (result error) {
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

	factory := NewFactory(formats, false, log)
	defer func() {
		if err := factory.Close(); err != nil && result == nil {
			result = err
		}
	}()

	writer, err := factory.Create(0)
	if err != nil {
		return fmt.Errorf("failed to setup output writer: %v", err)
	}

	if _, err := io.Copy(writer, reader); nil != err {
		result = fmt.Errorf("failed to process data: %v", err)
	}

	return
}
