package meta

import (
	"github.com/stretchr/testify/mock"
)

// InputMock implements Input as a mock for testing.
type InputMock struct {
	mock.Mock
}

// NewDummyInputMock returns an InputMock instance which returns dummy values.
func NewDummyInputMock() *InputMock {
	i := new(InputMock)
	i.On("Stem").Maybe().Return("")
	i.On("Size").Maybe().Return(int64(-1))

	return i
}

// NewInputMock returns an InputMock instance which describes the input with the given stem and an unknown size.
func NewInputMock(stem string) *InputMock {
	i := new(InputMock)
	i.On("Stem").Maybe().Return(stem)
	i.On("Size").Maybe().Return(int64(-1))

	return i
}

// Stem returns the mocked stem value.
func (m *InputMock) Stem() string {
	args := m.Called()
	return args.String(0)
}

// Size returns the mocked size value.
func (m *InputMock) Size() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}
