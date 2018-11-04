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

type Hash struct {
	algs   map[string]hash.Hash
	closed bool
}

func NewHash() *Hash {
	return &Hash{
		algs: map[string]hash.Hash{
			"md5":       md5.New(),
			"sha1":      sha1.New(),
			"sha256":    sha256.New(),
			"sha512":    sha512.New(),
			"ripemd160": ripemd160.New(),
		},
	}
}

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

func (h *Hash) Close() error {
	h.closed = true
	return nil
}

func (h *Hash) Md5() string {
	return h.Sum("md5")
}

func (h *Hash) Sha1() string {
	return h.Sum("sha1")
}

func (h *Hash) Sha256() string {
	return h.Sum("sha256")
}

func (h *Hash) Sha512() string {
	return h.Sum("sha512")
}

func (h *Hash) Ripemd160() string {
	return h.Sum("ripemd160")
}

func (h *Hash) Sum(a string) string {
	if !h.closed {
		return "undefined"
	}

	return fmt.Sprintf("%x", h.algs[a].Sum(nil))
}
