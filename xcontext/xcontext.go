package xcontext

import (
	"context"
	"time"
)

// WithFirstCancel creates a new context that gets canceled when ctx1 or ctx2 is canceled.
func WithFirstCancel(ctx1, ctx2 context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx1)
	stopf := context.AfterFunc(ctx2, func() {
		cancel()
	})
	return ctx, func() {
		cancel()
		stopf()
	}
}

// Detach returns a context that keeps all the values of its parent context
// but detaches from the cancellation and error handling.
func Detach(ctx context.Context) context.Context { return detachedContext{ctx} }

type detachedContext struct{ parent context.Context }

func (v detachedContext) Deadline() (time.Time, bool) { return time.Time{}, false }
func (v detachedContext) Done() <-chan struct{}       { return nil }
func (v detachedContext) Err() error                  { return nil }
func (v detachedContext) Value(key any) any           { return v.parent.Value(key) }
