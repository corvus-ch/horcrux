package restore

import (
	"io"
	"os"

	"github.com/howeyc/gopass"
)

type Prompter interface {
	Prompt(prompt string) ([]byte, error)
}

type prompter struct {
	r gopass.FdReader
	w io.Writer
}

func NewPrompter(w io.Writer) *prompter {
	return &prompter{os.Stdin, w}
}

func (p *prompter) Prompt(prompt string) ([]byte, error) {
	return gopass.GetPasswdPrompt(prompt, true, p.r, p.w)
}
