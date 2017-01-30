package actor_test

import (
	"sync"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func ExampleNewRootContext() {
	var wg sync.WaitGroup
	wg.Add(1)
	root := actor.NewRootContext()
	root.Meta().Add("Foo", "bar")

	props := actor.FromFunc(func(ctx actor.Context) {

	}).WithMiddleware(func(next actor.ActorFunc) actor.ActorFunc {
		fn := func(c actor.Context) {
			c.Meta().Add("parentid", "123")

			next(c)
		}

		return fn
	})

	actor.Spawn(props)
	wg.Wait()
}

// contextual := &ContextualEnvelope{
// 				Contextual: c.Meta("contextaul"),
// 				Message:    c.Message(),
// 			}
