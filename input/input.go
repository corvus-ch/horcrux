package input

// Input describes an interface used by formats to get information about the input.
type Input interface {
	// Name returns the base name of the input file
	Name() string

	// Path returns the absolute path of the input file
	Path() string

	// Stem returns the inputs file name without the extension.
	// It will be used to create the file names of the output.
	Stem() string

	// Size returns the size of the input in bytes.
	// If the input size can not be determined, a negative number will be returned.
	Size() int64
}
