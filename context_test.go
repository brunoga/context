package context

import (
	"testing"
	"time"
)

func TestWait_OneChild(t *testing.T) {
	root := Background()

	value := 0

	go func(ctx Context) {
		time.Sleep(1 * time.Millisecond)
		value = 1
		ctx.Finished()
	}(Child(root))

	root.Wait()

	if value != 1 {
		t.Errorf("Expected value to be 1. Got %d.", value)
	}
}

func TestWait_MultipleChildren(t *testing.T) {
	root := Background()

	value := 0

	go func(ctx Context) {
		time.Sleep(1 * time.Millisecond)
		value = 1
		ctx.Finished()
	}(Child(root))

	go func(ctx Context) {
		time.Sleep(2 * time.Millisecond)
		value = 2
		ctx.Finished()
	}(Child(root))

	go func(ctx Context) {
		time.Sleep(3 * time.Millisecond)
		value = 3
		ctx.Finished()
	}(Child(root))

	root.Wait()

	if value != 3 {
		t.Errorf("Expected value to be 3. Got %d.", value)
	}
}
