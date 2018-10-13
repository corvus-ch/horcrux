package restore

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bketelsen/logr"
	"github.com/corvus-ch/shamir"
	"golang.org/x/crypto/openpgp"
)

func Restore(cfg Config, p Prompter, log logr.Logger) error {
	var closers []io.Closer
	defer func() {
		for i := len(closers) - 1; i >= 0; i-- {
			closers[i].Close()
		}
	}()

	fileNames := cfg.FileNames()
	inputs := make(map[byte]io.Reader, len(fileNames))

	for _, fileName := range fileNames {
		log.Infof("Reading file %s\n", fileName)
		i, err := strconv.ParseUint(strings.TrimLeft(filepath.Ext(fileName), "."), 10, 8)
		if nil != err {
			return fmt.Errorf("failed to get shamir index from file name: %v", err)
		}

		file, err := os.Open(fileName)
		if nil != err {
			return fmt.Errorf("failed to open %s: %v", fileName, err)
		}

		closers = append(closers, file)

		format, err := cfg.Format()
		if err != nil {
			return fmt.Errorf("failed to to detect input format: %v", err)
		}
		in, err := format.Reader(file)
		if err != nil {
			return fmt.Errorf("failed to create input reader: %v", err)
		}

		if cfg.Decrypt() {
			in, err = gedDecryptionReader(in, p)
			if err != nil {
				return fmt.Errorf("failed to create decryption reader: %v", err)
			}
		}

		inputs[byte(i)] = in
	}

	in, err := shamir.NewReader(inputs)
	if nil != err {
		return fmt.Errorf("failed to create processing pipeline: %v", err)
	}

	out, err := cfg.Output()
	if nil != err {
		return fmt.Errorf("failed to open output: %v", err)
	}

	if c, ok := out.(io.Closer); ok {
		closers = append(closers, c)
	}

	if _, err := io.Copy(out, in); nil != err {
		return fmt.Errorf("Failed to process data: %v", err)
	}

	return nil
}

func gedDecryptionReader(r io.Reader, p Prompter) (io.Reader, error) {
	md, err := openpgp.ReadMessage(r, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return p.Prompt("Enter password: ")
	}, nil)
	if nil != err {
		return nil, err
	}

	return md.UnverifiedBody, nil
}
