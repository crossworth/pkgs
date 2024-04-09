package xtest

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"
)

// TestingT is an interface that we accept
// on functions related to tests and benchmarks.
type TestingT interface {
	Errorf(format string, args ...any)
	FailNow()
	Cleanup(fn func())
}

// Wait is a custom helper function used to wait for operations inside tests.
func Wait(t TestingT, timeout time.Duration, check func(t TestingT) bool) {
	if h, ok := t.(interface {
		Helper()
	}); ok {
		h.Helper()
	}
	result := make(chan bool)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	go func() {
		for {
			if check(t) {
				result <- true
				return
			}
			// timeout happened
			if ctx.Err() != nil {
				return
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()
	select {
	case <-ctx.Done():
		t.Errorf("timeout waiting: %s", timeout.String())
		t.FailNow()
	case <-result:
	}
}

// Must returns the value asserting that there are no errors.
func Must[T any](t TestingT) func(v T, err error) T {
	if h, ok := t.(interface {
		Helper()
	}); ok {
		h.Helper()
	}
	return func(v T, err error) T {
		require.NoError(t, err)
		return v
	}
}
