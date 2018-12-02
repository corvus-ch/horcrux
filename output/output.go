package output

// Output represents the list of output files grouped by output format.
// The definitive list of output files is unknown when this object gets created. The lists are therefore represented as
// channels. New files will get added to the channel as they are created and the channel will be closed once no more
// files have to be expected.
type Output map[string]chan File

// NewOutput creates a new instance of Output.
func NewOutput() Output {
	return Output{}
}

// Format returns the file list for the given output format.
func (o Output) Format(format string) chan File {
	ch, ok := o[format]
	if !ok {
		ch = make(chan File, 32)
		o[format] = ch
	}

	return ch
}

// Append adds a new output file for the given format.
func (o Output) Append(format, path string, meta map[string]interface{}) chan File {
	ch := o.Format(format)
	ch <- NewFile(path, meta)
	return ch
}
