package floats

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNearlyEqual(t *testing.T) {
	t.Parallel()
	require.True(t, NearlyEqual(10.10, 10.10, 0))
	require.True(t, NearlyEqual(10.10001, 10.10, 0.1))
	require.True(t, NearlyEqual(10.10011, 10.10, 0.1))
	require.True(t, NearlyEqual(10.10111, 10.10, 0.1))
	require.True(t, NearlyEqual(10.11111, 10.10, 0.1))
	require.True(t, NearlyEqual(10.0+0.1, 10.1, 0))
	require.False(t, NearlyEqual(10.0+0.12, 10.1, 0))
	require.False(t, NearlyEqual(10.0+0.12, 10.1, 0.019))
	require.True(t, NearlyEqual(10.0+0.12, 10.1, 0.020))
}
