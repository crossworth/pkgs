package xtime

import (
	"context"
	"time"
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
