package output

type Output map[string]chan File

func NewOutput() Output {
	return Output{}
}

func (o Output) Format(format string) chan File {
	ch, ok := o[format]
	if !ok {
		ch = make(chan File, 32)
		o[format] = ch
	}

	return ch
}

func (o Output) Append(format, path string, meta map[string]interface{}) chan File {
	ch := o.Format(format)
	ch <- NewFile(path, meta)
	return ch
}
