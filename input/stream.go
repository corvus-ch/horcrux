package input

func NewStreamInput(stem string) Input {
	i := &stream{"part"}
	if len(stem) > 0 {
		i.stem = stem
	}

	return i
}

type stream struct {
	stem string
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
