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
	parentWg() *sync.WaitGroup
	childrenWg() *sync.WaitGroup
}

var Canceled = context.Canceled
var DeadlineExceeded = context.DeadlineExceeded

type ctxImpl struct {
	context.Context

	pWg *sync.WaitGroup
	cWg sync.WaitGroup
}

func (c *ctxImpl) Finished() {
	c.pWg.Done()
}

func (c *ctxImpl) Wait() {
	c.cWg.Wait()
}

func (c *ctxImpl) context() context.Context {
	return c.Context
}

func (c *ctxImpl) parentWg() *sync.WaitGroup {
	return c.pWg
}

func (c *ctxImpl) childrenWg() *sync.WaitGroup {
	return &c.cWg
}

type emptyCtx ctxImpl

func (e *emptyCtx) Finished() {
	return
}

func (e *emptyCtx) Wait() {
	e.cWg.Wait()
}

func (e *emptyCtx) context() context.Context {
	return e.Context
}

func (c *emptyCtx) parentWg() *sync.WaitGroup {
	return nil
}

func (e *emptyCtx) childrenWg() *sync.WaitGroup {
	return &e.cWg
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
	parent.parentWg().Add(1)
	ctx, c := context.WithCancel(parent.context())
	return &ctxImpl{
		ctx,
		parent.childrenWg(),
		sync.WaitGroup{},
	}, CancelFunc(c)
}

func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc) {
	parent.parentWg().Add(1)
	ctx, c := context.WithDeadline(parent.context(), deadline)
	return &ctxImpl{
		ctx,
		parent.childrenWg(),
		sync.WaitGroup{},
	}, CancelFunc(c)
}

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	parent.parentWg().Add(1)
	ctx, c := context.WithTimeout(parent.context(), timeout)
	return &ctxImpl{
		ctx,
		parent.childrenWg(),
		sync.WaitGroup{},
	}, CancelFunc(c)
}

func WithContext(parent context.Context) Context {
	return &ctxImpl{
		parent,
		nil,
		sync.WaitGroup{},
	}
}
