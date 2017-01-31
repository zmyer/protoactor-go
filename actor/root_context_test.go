package actor

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewRootContext(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	root := NewRootContext()
	root.Meta().Add("Foo", "bar")

	props := FromFunc(func(ctx Context) {
		if ctx.Message() == "hello" {
			meta := ctx.Meta().Get("contextual")
			fmt.Println(meta)
			wg.Done()
		}

	}).WithMiddleware(func(next ActorFunc) ActorFunc {
		fn := func(c Context) {
			c.Meta().Add("contextual", "hello world")
			next(c)
		}
		return fn
	})

	pid := Spawn(props)
	root.Tell("hello", pid)
	wg.Wait()
}
