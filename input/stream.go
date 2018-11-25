package input

import "github.com/corvus-ch/horcrux/hash"

// NewStreamInput returns an instance of Input representing an input passed via STDIN.
func NewStreamInput(stem string) Input {
	i := &stream{
		stem: "part",
		hash: hash.NewHash(),
	}
	if len(stem) > 0 {
		i.stem = stem
	}

	return i
}

type stream struct {
	stem string
	hash *hash.Hash
}

// Name will return an empty string.
func (i *stream) Name() string {
	return ""
}

// Path will return an empty string.
func (i *stream) Path() string {
	return ""
}

// Path returns the inputs stem.
// By default, this will be 'part' or whatever was set by the output
// option.
func (i *stream) Stem() string {
	return i.stem
}

// Size will be -1.
func (i *stream) Size() int64 {
	return -1
}

// Checksum returns the inputs checksum calculated for the given algorithm.
// The checksum will not be valid until the checksum object got closed.
func (i *stream) Checksum(alg string) (string, error) {
	return i.hash.Sum(alg)
}

// Write will add the bytes to the checksum calculation.
func (i *stream) Write(p []byte) (int, error) {
	return i.hash.Write(p)
}

// Close will closes the checksum calculation making them valid.
func (i *stream) Close() error {
	return i.hash.Close()
}
