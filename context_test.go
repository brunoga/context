package context

import (
	"testing"
	"time"
)

func TestWait_OneChild(t *testing.T) {
	parent := Background()

	ctx, cancel := WithCancel(parent)
	defer cancel()

	value := 0

	go func(ctx Context) {
		time.Sleep(1 * time.Millisecond)
		value = 1
		ctx.Finished()
	}(EnableWait(ctx))

	parent.WaitForChildren()

	if value != 1 {
		t.Errorf("Expected value to be 1. Got %d.", value)
	}
}

func TestWait_MultipleChildren(t *testing.T) {
	parent := Background()

	ctx, cancel := WithCancel(parent)
	defer cancel()

	value := 0

	go func(ctx Context) {
		time.Sleep(1 * time.Millisecond)
		value = 1
		ctx.Finished()
	}(EnableWait(ctx))

	go func(ctx Context) {
		time.Sleep(2 * time.Millisecond)
		value = 2
		ctx.Finished()
	}(EnableWait(ctx))

	go func(ctx Context) {
		time.Sleep(3 * time.Millisecond)
		value = 3
		ctx.Finished()
	}(EnableWait(ctx))

	parent.WaitForChildren()

	if value != 3 {
		t.Errorf("Expected value to be 3. Got %d.", value)
	}
}
