package input

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"

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
// Checksums will not be considered valid until Close() got
// called.
func (h *Hash) Close() error {
	h.closed = true
	return nil
}

// Md5 returns the inputs MD5 checksums.
func (h *Hash) Md5() string {
	return h.Sum(MD5)
}

// Sha1 returns the inputs SHA1 checksums.
func (h *Hash) Sha1() string {
	return h.Sum(SHA1)
}

// Sha256 returns the inputs SHA256 checksums.
func (h *Hash) Sha256() string {
	return h.Sum(SHA256)
}

// Sha512 returns the inputs SHA512 checksums.
func (h *Hash) Sha512() string {
	return h.Sum(SHA512)
}

// Ripemd160 returns the inputs RIPEMD160 checksums.
func (h *Hash) Ripemd160() string {
	return h.Sum(RIPEMD160)
}

// Sum returns the inputs checksums for the given algorithm.
func (h *Hash) Sum(a string) string {
	if !h.closed {
		return "undefined"
	}

	return fmt.Sprintf("%x", h.algs[a].Sum(nil))
}
