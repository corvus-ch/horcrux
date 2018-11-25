package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"os"

	"golang.org/x/crypto/ripemd160"
)

const (
	// MD5 holds the identifier of the MD5 algorithm.
	MD5 = "md5"

	// SHA1 holds the identifier of the SHA1 algorithm.
	SHA1 = "sha1"

	// SHA256 holds the identifier of the SHA256 algorithm.
	SHA256 = "sha256"

	// SHA512 holds the identifier of the SHA512 algorithm.
	SHA512 = "sha512"

	// RIPEMD160 holds the identifier of the RIPEMD160 algorithm.
	RIPEMD160 = "ripemd160"
)

// Hash represents a set of input checksums calculated with several algorithms.
type Hash struct {
	algs   map[string]hash.Hash
	closed bool
}

// NewHash returns a new instance of Hash with pre populated algorithms.
func NewHash() *Hash {
	return &Hash{
		algs: map[string]hash.Hash{
			MD5:       md5.New(),
			SHA1:      sha1.New(),
			SHA256:    sha256.New(),
			SHA512:    sha512.New(),
			RIPEMD160: ripemd160.New(),
		},
	}
}

func NewHashForPath(path string) (cs *Hash, err error) {
	cs = NewHash()
	f, _ := os.Open(path)
	defer f.Close()
	if f != nil {
		if _, err = io.Copy(cs, f); err != nil {
			return
		}
	}

	if err = cs.Close(); err != nil {
		return
	}

	return
}

// Write implements io.Writer and updates the all the checksums.
func (h *Hash) Write(p []byte) (n int, err error) {
	for _, w := range h.algs {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			return
		}
	}
	return len(p), nil
}

// Close implements io.Closer.
// Checksums will not be considered valid until Close() got called.
func (h *Hash) Close() error {
	h.closed = true
	return nil
}

// Sum returns the inputs checksum for the given algorithm.
func (h *Hash) Sum(a string) (string, error) {
	if !h.closed {
		return "undefined", nil
	}

	alg, ok := h.algs[a]
	if !ok {
		return "", fmt.Errorf("unknown checksum algorithm %v", a)
	}

	return fmt.Sprintf("%x", alg.Sum(nil)), nil
}
