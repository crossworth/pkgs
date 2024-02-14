package th

import (
	"errors"
	"fmt"
)

var _ interface {
	TestingT
	Helper()
	Cleanup(func())
} = (*FakeT)(nil)

// FakeT is a fake testing input
// that can be used in places where the testing.T is not possible.
type FakeT struct {
	message string
	failNow bool
	cleanup []func()
}

func (s *FakeT) Errorf(format string, args ...any) {
	s.message = fmt.Sprintf(format, args...)
}

func (s *FakeT) FailNow() {
	s.failNow = true
}

func (s *FakeT) Helper() {
}

func (s *FakeT) Cleanup(f func()) {
	s.cleanup = append(s.cleanup, f)
}

func (s *FakeT) Check() error {
	if s.failNow {
		for i := len(s.cleanup); i > 0; i-- {
			if s.cleanup != nil {
				s.cleanup[i]()
			}
		}
		return errors.New(s.message)
	}
	return nil
}
