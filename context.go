package context

import (
	"context"
	"sync"
	"time"
)

type Context interface {
	context.Context

	Finished()
	Wait()

	context() context.Context
	wg() *sync.WaitGroup
}

var Canceled = context.Canceled
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
