// Package context implements a drop-in replacement for the standard library
// context package. It has all the features of the standard library version plus
// support for waiting on derived contexts (including supporting
// functions/methods).
//
// See https://golang.org/pkg/context.
//
// Example usage:
//
// root := context.Background()
//
// workersCtx, cancel := context.WithCancel(root)
// defer cancel()
//
// for i:=0; i < numWorkers; i++ {
//	   go startWorker(context.EnableWait(workersCtx))
// }
//
// root.WaitForChildren()
package context

import (
	"context"
	"sync"
	"time"
)

var (
	// Errors.
	Canceled         = context.Canceled
	DeadlineExceeded = context.DeadlineExceeded
)

// Context behaves exactly like a standard library Context but also includes
// support for waiting on derived (child) Contexts.
//
// See https://golang.org/pkg/context/#Context.
type Context interface {
	context.Context

	// Finished reports back to the parent Context that the work associated
	// with this Context has finished. This must be explicitly called when
	// using the waiting feature.
	Finished()

	// Wait waits on all immediate children to finish their work. It blocks
	// until all children report that their work is finished.
	WaitForChildren()

	context() context.Context

	pWg() *sync.WaitGroup
	cWg() *sync.WaitGroup
}

type ctxImpl struct {
	context.Context

	parentWg   *sync.WaitGroup
	childrenWg sync.WaitGroup
}

func (c *ctxImpl) Finished() {
	if c.parentWg != nil {
		// Only non-root contexts have parents.
		c.parentWg.Done()
	}
}

func (c *ctxImpl) WaitForChildren() {
	c.childrenWg.Wait()
}

func (c *ctxImpl) context() context.Context {
	return c.Context
}

func (c *ctxImpl) pWg() *sync.WaitGroup {
	return c.parentWg
}

func (c *ctxImpl) cWg() *sync.WaitGroup {
	return &c.childrenWg
}

func (c *ctxImpl) Err() error {
	switch c.Context.Err() {
	case context.Canceled:
		return Canceled
	case context.DeadlineExceeded:
		return DeadlineExceeded
	}

	// Got an unexpected error. Just return it.
	//
	// TODO(bga): Maybe log something here?
	return c.Context.Err()
}

func Background() Context {
	return &ctxImpl{
		context.Background(),
		nil,
		sync.WaitGroup{},
	}
}

func TODO() Context {
	return &ctxImpl{
		context.TODO(),
		nil,
		sync.WaitGroup{},
	}
}

type CancelFunc context.CancelFunc

func WithCancel(parent Context) (Context, CancelFunc) {
	ctx, c := context.WithCancel(parent.context())
	return &ctxImpl{
		ctx,
		parent.cWg(),
		sync.WaitGroup{},
	}, CancelFunc(c)
}

func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc) {
	ctx, c := context.WithDeadline(parent.context(), deadline)
	return &ctxImpl{
		ctx,
		parent.cWg(),
		sync.WaitGroup{},
	}, CancelFunc(c)
}

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	ctx, c := context.WithTimeout(parent.context(), timeout)
	return &ctxImpl{
		ctx,
		parent.cWg(),
		sync.WaitGroup{},
	}, CancelFunc(c)
}

// EnableWait enables waiting on this context completion. When the work
// associated with this context finishes (ctx.Finished() is called the same
// number of times that EnableWait() is called), any caller waiting on the
// parent context will unblock.
func EnableWait(ctx Context) Context {
	ctx.pWg().Add(1)

	return ctx
}
