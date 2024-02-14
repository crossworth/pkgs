package xtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWithinDuration(t *testing.T) {
	t.Parallel()
	require.True(t, WithinDuration(time.Now(), time.Now().Add(1*time.Second), 2*time.Second))
	require.False(t, WithinDuration(time.Now(), time.Now().Add(5*time.Second), 1*time.Second))
}

func TestDate(t *testing.T) {
	t.Parallel()
	d := Date(2021, 12, 22)
	require.Equal(t, 0, d.Hour())
	require.Equal(t, 0, d.Minute())
	require.Equal(t, 0, d.Second())
	require.Equal(t, 0, d.Nanosecond())
}

func TestStartOfMonth(t *testing.T) {
	t.Parallel()
	d := time.Date(2021, 12, 22, 19, 40, 10, 0, time.UTC)
	require.Equal(t, "2021-12-01 00:00:00 +0000 UTC", StartOfMonth(d).String())

	d = time.Date(2021, 12, 1, 19, 40, 10, 0, time.UTC)
	require.Equal(t, "2021-12-01 00:00:00 +0000 UTC", StartOfMonth(d).String())

	d = time.Date(2021, 12, 31, 19, 40, 10, 0, time.UTC)
	require.Equal(t, "2021-12-01 00:00:00 +0000 UTC", StartOfMonth(d).String())

	d = time.Date(2021, 12, 31, 23, 59, 59, 59, time.UTC)
	require.Equal(t, "2021-12-01 00:00:00 +0000 UTC", StartOfMonth(d).String())
}

func TestEndOfMonth(t *testing.T) {
	t.Parallel()
	d := time.Date(2021, 12, 22, 19, 40, 10, 0, time.UTC)
	require.Equal(t, "2021-12-31 23:59:59.999999999 +0000 UTC", EndOfMonth(d).String())

	d = time.Date(2021, 12, 1, 19, 40, 10, 0, time.UTC)
	require.Equal(t, "2021-12-31 23:59:59.999999999 +0000 UTC", EndOfMonth(d).String())

	d = time.Date(2021, 12, 31, 19, 40, 10, 0, time.UTC)
	require.Equal(t, "2021-12-31 23:59:59.999999999 +0000 UTC", EndOfMonth(d).String())

	d = time.Date(2021, 12, 31, 23, 59, 59, 59, time.UTC)
	require.Equal(t, "2021-12-31 23:59:59.999999999 +0000 UTC", EndOfMonth(d).String())

	d = time.Date(2021, 2, 1, 1, 1, 1, 1, time.UTC)
	require.Equal(t, "2021-02-28 23:59:59.999999999 +0000 UTC", EndOfMonth(d).String())
}

func TestStartOfDay(t *testing.T) {
	t.Parallel()
	d := time.Date(2021, 12, 22, 19, 40, 10, 0, time.UTC)
	require.Equal(t, "2021-12-22 00:00:00 +0000 UTC", StartOfDay(d).String())
	d = time.Date(2021, 12, 3, 19, 40, 10, 0, time.UTC)
	require.Equal(t, "2021-12-03 00:00:00 +0000 UTC", StartOfDay(d).String())
	d = time.Date(2021, 12, 31, 19, 40, 10, 0, time.UTC)
	require.Equal(t, "2021-12-31 00:00:00 +0000 UTC", StartOfDay(d).String())
	d = time.Date(2021, 12, 20, 23, 59, 59, 59, time.UTC)
	require.Equal(t, "2021-12-20 00:00:00 +0000 UTC", StartOfDay(d).String())
}

func TestEndOfDay(t *testing.T) {
	t.Parallel()
	d := time.Date(2021, 12, 22, 19, 40, 10, 0, time.UTC)
	require.Equal(t, "2021-12-22 23:59:59.999999999 +0000 UTC", EndOfDay(d).String())
	d = time.Date(2021, 12, 1, 19, 40, 10, 0, time.UTC)
	require.Equal(t, "2021-12-01 23:59:59.999999999 +0000 UTC", EndOfDay(d).String())
	d = time.Date(2021, 12, 31, 19, 40, 10, 0, time.UTC)
	require.Equal(t, "2021-12-31 23:59:59.999999999 +0000 UTC", EndOfDay(d).String())
	d = time.Date(2021, 12, 31, 23, 59, 59, 59, time.UTC)
	require.Equal(t, "2021-12-31 23:59:59.999999999 +0000 UTC", EndOfDay(d).String())
	d = time.Date(2021, 2, 1, 1, 1, 1, 1, time.UTC)
	require.Equal(t, "2021-02-01 23:59:59.999999999 +0000 UTC", EndOfDay(d).String())
}

func TestAddMonths(t *testing.T) {
	t.Parallel()
	d := time.Date(2021, 3, 30, 19, 19, 19, 999999999, time.UTC)
	n := AddMonths(d, -1)
	require.Equal(t, "2021-02-28 19:19:19.999999999 +0000 UTC", n.String())
}
