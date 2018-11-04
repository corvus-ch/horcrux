package input

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

func (i *stream) Name() string {
	return ""
}

func (i *stream) Path() string {
	return ""
}

func (i *stream) Stem() string {
	return i.stem
}

func (i *stream) Size() int64 {
	return -1
}

func (i *stream) Checksums() *Hash {
	return i.checksumms
}

func (i *stream) Write(p []byte) (int, error) {
	return i.checksumms.Write(p)
}

func (i *stream) Close() error {
	return i.checksumms.Close()
}
