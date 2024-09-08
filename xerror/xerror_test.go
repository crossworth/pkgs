package xerror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsErrNotFound(t *testing.T) {
	t.Parallel()
	require.True(t, IsErrNotFound(MakeNotFoundError()))
}

func TestEncodeSliceErrParams(t *testing.T) {
	t.Parallel()
	err := MakeBadRequestError(
		ErrParam("reason", "invalid argument"),
		ErrParam("valueProvided", "-1"),
		ErrParam("parseError", errors.New("parse error")),
	)
	require.Equal(t, "bad_request: parseError=parse error, reason=invalid argument, valueProvided=-1", err.Error())

	err = MakeBadRequestError(
		ErrParam("field", map[string]any{
			"reason":        "invalid argument",
			"valueProvided": "-1",
			"parseError":    errors.New("parse error"),
		}),
	)
	require.Equal(t, "bad_request: field=map[parseError:parse error reason:invalid argument valueProvided:-1]", err.Error())
}
