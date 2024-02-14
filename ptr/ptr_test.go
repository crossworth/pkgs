package ptr

import (
	"fmt"
	"testing"
	"time"

	"github.com/crossworth/pkgs/floats"
	"github.com/stretchr/testify/require"
)

func TestOf(t *testing.T) {
	t.Parallel()
	i := 10
	require.Equal(t, &i, Of(i))
	f := 10.0
	require.Equal(t, &f, Of(f))
	tm := time.Now()
	require.Equal(t, &tm, Of(tm))
}

func TestValueOrDefault(t *testing.T) {
	t.Parallel()
	require.Equal(t, 10, ValueOrDefault(nil, 10))
	i := 5
	require.Equal(t, 5, ValueOrDefault(&i, 10))
}

func TestValuePtrOrNil(t *testing.T) {
	t.Parallel()
	require.Equal(t, (*int)(nil), ValuePtrOrNil(10, true))
	require.Equal(t, Of(10), ValuePtrOrNil(10, false))
}

func TestValueOf(t *testing.T) {
	t.Parallel()
	require.Equal(t, 10, ValueOf(Of(10)))
}

func TestFloat32To64(t *testing.T) {
	t.Parallel()
	v := float32(10.30)
	r := *Float32To64(&v)
	require.True(t, floats.NearlyEqual(r, 10.30, 0.001))
	require.Equal(t, "float64", fmt.Sprintf("%T", r))
}
