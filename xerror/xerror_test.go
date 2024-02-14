package xerror

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsErrNotFound(t *testing.T) {
	t.Parallel()
	require.True(t, IsErrNotFound(MakeNotFoundError()))
}
