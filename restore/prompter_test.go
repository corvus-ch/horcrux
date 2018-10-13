package restore_test

import (
	"github.com/pkg/errors"
)

type Prompter struct {
	passwords []string
}

func NewPrompter(passwords []string) *Prompter {
	return &Prompter{passwords: passwords}
}

func (p *Prompter) Prompt(_ string) ([]byte, error) {
	if len(p.passwords) < 1 {
		return nil, errors.New("no passwords available")
	}
	password := p.passwords[0]
	p.passwords = p.passwords[1:]

	return []byte(password), nil
}
