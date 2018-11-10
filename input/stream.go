package input

// NewStreamInput returns an instance of Input representing an input passed via STDIN.
func NewStreamInput(stem string) Input {
	i := &stream{
		stem:       "part",
		checksumms: NewHash(),
	}
	if len(stem) > 0 {
		i.stem = stem
	}

	return i
}

type stream struct {
	stem       string
	checksumms *Hash
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

// Checksums returns a set containing the inputs checksums calculated
// with several algorithms. The checksum will not be valid until the
// checksum object will be closed.
func (i *stream) Checksums() *Hash {
	return i.checksumms
}

// Write will add the bytes to the checksums.
func (i *stream) Write(p []byte) (int, error) {
	return i.checksumms.Write(p)
}

// Close will closes the checksum calculation making them valid.
func (i *stream) Close() error {
	return i.checksumms.Close()
}
