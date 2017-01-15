package actor

import (
	"testing"

	"github.com/AsynkronIT/protoactor-go/process"
	"github.com/stretchr/testify/assert"
)

func TestDeadLetterAfterStop(t *testing.T) {
	actor := Spawn(FromProducer(NewBlackHoleActor))
	done := false
	sub := process.EventStream.Subscribe(func(msg interface{}) {
		if deadLetter, ok := msg.(*process.DeadLetter); ok {
			if deadLetter.PID == actor {
				done = true
			}
		}
	})
	defer process.EventStream.Unsubscribe(sub)

	stopActorFuture(actor).Wait()

	actor.Tell("hello")

	assert.True(t, done)
}
