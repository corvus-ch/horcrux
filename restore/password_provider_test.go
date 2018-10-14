package restore_test

import (
	"github.com/pkg/errors"
)

type TestPasswordProvider struct {
	passwords []string
}

func NewPasswordProvider(passwords []string) *TestPasswordProvider {
	return &TestPasswordProvider{passwords: passwords}
}

func (p *TestPasswordProvider) GetPassword(_ string) ([]byte, error) {
	if len(p.passwords) < 1 {
		return nil, errors.New("no passwords available")
	}
	password := p.passwords[0]
	p.passwords = p.passwords[1:]

	return []byte(password), nil
}
