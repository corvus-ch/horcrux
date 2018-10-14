package restore

import (
	"io"
	"os"

	"github.com/howeyc/gopass"
)

// PasswordProvider allows Restore to ask for passwords of encrypted parts.
type PasswordProvider interface {
	GetPassword(prompt string) ([]byte, error)
}

// StdinPasswordProvider implements PasswordProvider
// It uses a given writer for writing input prompts and STDIN for reading passwords from.
type StdinPasswordProvider struct {
	r gopass.FdReader
	w io.Writer
}

// NewPasswordProvider returns an instance of StdinPasswordProvider.
func NewPasswordProvider(w io.Writer) *StdinPasswordProvider {
	return &StdinPasswordProvider{os.Stdin, w}
}

// GetPassword prompts the user and returns the password read from the terminal.
func (p *StdinPasswordProvider) GetPassword(prompt string) ([]byte, error) {
	return gopass.GetPasswdPrompt(prompt, true, p.r, p.w)
}
