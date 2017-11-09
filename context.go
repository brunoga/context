// Package context implements a drop-in replacement for the standard library
// context package. It has all the features of the standard library version plus
// support for waiting on derived contexts (including supporting
// functions/methods).
//
// See https://golang.org/pkg/context.
package context

import (
	"context"
	"sync"
	"time"
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
	Wait()

	context() context.Context
	wg() *sync.WaitGroup
}

// Canceled just exposes the standard library context.Canceled error.
//
// See https://golang.org/pkg/context/#pkg-variables.
var Canceled = context.Canceled

// DeadlineExceeded just exposes the standard library context.DeadlineExceeded
// error.
//
// See https://golang.org/pkg/context/#pkg-variables.
var DeadlineExceeded = context.DeadlineExceeded

type ctxImpl struct {
	context.Context

	parentWg   *sync.WaitGroup
	childrenWg sync.WaitGroup
}

func (c *ctxImpl) Finished() {
	c.parentWg.Done()
}

func (c *ctxImpl) Wait() {
	c.childrenWg.Wait()
}

func (c *ctxImpl) context() context.Context {
	return c.Context
}

func (c *ctxImpl) wg() *sync.WaitGroup {
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

type emptyCtx ctxImpl

func (e *emptyCtx) Finished() {
	return
}

func (e *emptyCtx) Wait() {
	e.childrenWg.Wait()
}

func (e *emptyCtx) context() context.Context {
	return e.Context
}

func (e *emptyCtx) wg() *sync.WaitGroup {
	return &e.childrenWg
}

func Background() Context {
	return &emptyCtx{
		context.Background(),
		nil,
		sync.WaitGroup{},
	}
}

func TODO() Context {
	return &emptyCtx{
		context.TODO(),
		nil,
		sync.WaitGroup{},
	}
}

type CancelFunc context.CancelFunc

func WithCancel(parent Context) (Context, CancelFunc) {
	parent.wg().Add(1)
	ctx, c := context.WithCancel(parent.context())
	return &ctxImpl{
		ctx,
		parent.wg(),
		sync.WaitGroup{},
	}, CancelFunc(c)
}

func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc) {
	parent.wg().Add(1)
	ctx, c := context.WithDeadline(parent.context(), deadline)
	return &ctxImpl{
		ctx,
		parent.wg(),
		sync.WaitGroup{},
	}, CancelFunc(c)
}

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	parent.wg().Add(1)
	ctx, c := context.WithTimeout(parent.context(), timeout)
	return &ctxImpl{
		ctx,
		parent.wg(),
		sync.WaitGroup{},
	}, CancelFunc(c)
}

// WithStandardContext returns the parent standard library context.Context
// Wrapped in a Cntext instance. The returned Context can then make use of the
// waiting functionality.
//
// Everything set in the parent Context (deadlines, cancel function, values)
// will be present in the returne one.
func WithStandardContext(parent context.Context) Context {
	return &ctxImpl{
		parent,
		nil,
		sync.WaitGroup{},
	}
}

// WithContext returns the parent context with its internal wait group counter
// incremented. This can be used to simply pass the same context to different
// goroputines without having to create actual new Context instances.
//
// For general Context purposes, the returned Context is, effectivelly, the
// parent context itself.
//
// For waiting purposes, the returned context can be considered a derived
// (child) one).
func WithContext(parent Context) Context {
	parent.wg().Add(1)
	return parent
}
