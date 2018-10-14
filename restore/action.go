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

// Restore merges a set of parts into the original secret.
func Restore(cfg Config, p PasswordProvider, log logr.Logger) error {
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
		i, in, file, err := getInput(cfg, fileName, p)
		if file != nil {
			closers = append(closers, file)
		}
		if err != nil {
			return err
		}
		inputs[byte(i)] = in
	}

	in, err := shamir.NewReader(inputs)
	if err != nil {
		return fmt.Errorf("failed to create processing pipeline: %v", err)
	}

	out, err := cfg.Output()
	if err != nil {
		return fmt.Errorf("failed to open output: %v", err)
	}

	if c, ok := out.(io.Closer); ok {
		closers = append(closers, c)
	}

	return restore(in, out)
}

func getInput(cfg Config, fileName string, p PasswordProvider) (uint64, io.Reader, io.ReadCloser, error) {
	i, err := strconv.ParseUint(strings.TrimLeft(filepath.Ext(fileName), "."), 10, 8)
	if err != nil {
		return i, nil, nil, fmt.Errorf("failed to get shamir index from file name: %v", err)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return i, nil, file, fmt.Errorf("failed to open %s: %v", fileName, err)
	}

	in, err := getFormatReader(cfg, file, p)

	return i, in, file, err
}

func getFormatReader(cfg Config, file io.Reader, p PasswordProvider) (io.Reader, error) {
	format, err := cfg.Format()
	if err != nil {
		return nil, fmt.Errorf("failed to to detect input format: %v", err)
	}

	in, err := format.Reader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create input reader: %v", err)
	}

	if cfg.Decrypt() {
		in, err = gedDecryptionReader(in, p)
		if err != nil {
			return nil, fmt.Errorf("failed to create decryption reader: %v", err)
		}
	}

	return in, nil
}

func gedDecryptionReader(r io.Reader, p PasswordProvider) (io.Reader, error) {
	md, err := openpgp.ReadMessage(r, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return p.GetPassword("Enter password: ")
	}, nil)
	if nil != err {
		return nil, err
	}

	return md.UnverifiedBody, nil
}

func restore(in io.Reader, out io.Writer) error {
	_, err := io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("failed to process data: %v", err)
	}

	return nil
}
