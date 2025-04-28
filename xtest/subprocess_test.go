package xtest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunInSubprocess(t *testing.T) {
	RunInSubprocess(t)
	t.Log("executed in sub process")
	require.True(t, true)
}
