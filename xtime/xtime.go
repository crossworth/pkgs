package xtime

import (
	"context"
	"time"

	"github.com/uniplaces/carbon"
)

// SleepWithContext sleeps respecting the given context.
func SleepWithContext(ctx context.Context, d time.Duration) {
	timer := time.NewTimer(d)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	case <-timer.C:
	}
}

// WithinDuration compares two dates and check if they are close using a delta.
func WithinDuration(actual time.Time, expected time.Time, delta time.Duration) bool {
	dt := expected.Sub(actual)
	return dt >= -delta && dt <= delta
}

// Date creates a date with hour, minute, second and nanosecond.
func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// StartOfMonth returns the beginning of the month.
func StartOfMonth(date time.Time) time.Time {
	y, m, _ := date.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, date.Location())
}

// EndOfMonth returns the end of the month.
func EndOfMonth(date time.Time) time.Time {
	return StartOfMonth(date).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// StartOfDay returns the beginning of the day.
func StartOfDay(date time.Time) time.Time {
	y, m, d := date.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, date.Location())
}

// EndOfDay returns the end of the day.
func EndOfDay(date time.Time) time.Time {
	return StartOfDay(date).AddDate(0, 0, 1).Add(-time.Nanosecond)
}

// AddMonths adds months from the time.Time provided
// respecting the calendar rules.
// That means that 31/03 -1 month will be 28/02.
func AddMonths(date time.Time, m int) time.Time {
	c := carbon.NewCarbon(date)
	return c.AddMonthsNoOverflow(m).Time
}
