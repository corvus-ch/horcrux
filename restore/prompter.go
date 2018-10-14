package restore

import (
	"io"
	"os"

	"github.com/howeyc/gopass"
)

// Prompter allows Restore to ask for passwords of encrypted parts.
type Prompter interface {
	Prompt(prompt string) ([]byte, error)
}

// StdinPrompter implements Prompter
// It uses a given writer for writing input prompts and STDIN for reading passwords from.
type StdinPrompter struct {
	r gopass.FdReader
	w io.Writer
}

// NewPrompter returns an instance of StdinPrompter.
func NewPrompter(w io.Writer) *StdinPrompter {
	return &StdinPrompter{os.Stdin, w}
}

func (p *StdinPrompter) Prompt(prompt string) ([]byte, error) {
	return gopass.GetPasswdPrompt(prompt, true, p.r, p.w)
}
