package th

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWait(t *testing.T) {
	t.Parallel()
	t.Run("wait will be executed", func(t *testing.T) {
		t.Parallel()
		nCalls := 0
		Wait(t, 1*time.Second, func(t TestingT) bool {
			nCalls++
			return true
		})
		require.Equal(t, 1, nCalls)
	})
	t.Run("timeout happens", func(t *testing.T) {
		t.Parallel()
		ctf := &captureTestFailure{}
		Wait(ctf, 1*time.Millisecond, func(t TestingT) bool {
			return false
		})
		require.Equal(t, "timeout waiting: %s", ctf.format)
		require.Equal(t, []any{"1ms"}, ctf.args)
		require.True(t, ctf.failNow)
	})
	t.Run("timeout happens (check slow)", func(t *testing.T) {
		t.Parallel()
		ctf := &captureTestFailure{}
		Wait(ctf, 1*time.Millisecond, func(t TestingT) bool {
			time.Sleep(1 * time.Second)
			return false
		})
		require.Equal(t, "timeout waiting: %s", ctf.format)
		require.Equal(t, []any{"1ms"}, ctf.args)
		require.True(t, ctf.failNow)
	})
}

type captureTestFailure struct {
	format  string
	args    []any
	failNow bool
	cleanup []func()
}

func (c *captureTestFailure) Errorf(format string, args ...any) {
	c.format = format
	c.args = args
}

func (c *captureTestFailure) FailNow() {
	c.failNow = true
}

func (c *captureTestFailure) Cleanup(fn func()) {
	c.cleanup = append(c.cleanup, fn)
}

func TestMust(t *testing.T) {
	require.Equal(t, 10, Must[int](t)(10, nil))
}
