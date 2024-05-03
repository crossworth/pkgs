package xtest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenPort(t *testing.T) {
	p := OpenPort(t)
	require.NotEmpty(t, p)
}
