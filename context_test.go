package context

import (
	"testing"
	"time"
)

func TestWait_OneChild(t *testing.T) {
	root := Background()

	child, cancel := WithCancel(root)
	defer cancel()

	value := 0

	go func(ctx Context) {
		time.Sleep(1 * time.Millisecond)
		value = 1
		ctx.Finished()
	}(child)

	root.Wait()

	if value != 1 {
		t.Errorf("Expected value to be 1. Got %d.", value)
	}
}

func TestWait_MultipleChildren(t *testing.T) {
	root := Background()

	child1, cancel1 := WithCancel(root)
	defer cancel1()

	child2, cancel2 := WithCancel(root)
	defer cancel2()

	child3, cancel3 := WithCancel(root)
	defer cancel3()

	value := 0

	go func(ctx Context) {
		time.Sleep(1 * time.Millisecond)
		value = 1
		ctx.Finished()
	}(child1)

	go func(ctx Context) {
		time.Sleep(2 * time.Millisecond)
		value = 2
		ctx.Finished()
	}(child2)

	go func(ctx Context) {
		time.Sleep(3 * time.Millisecond)
		value = 3
		ctx.Finished()
	}(child3)

	root.Wait()

	if value != 3 {
		t.Errorf("Expected value to be 3. Got %d.", value)
	}
}
