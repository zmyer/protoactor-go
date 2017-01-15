package actor

import (
	"testing"

	"github.com/AsynkronIT/protoactor-go/process"
	"github.com/stretchr/testify/assert"
)

type EchoOnStartActor struct{ replyTo *process.ID }

func (state *EchoOnStartActor) Receive(context Context) {
	switch context.Message().(type) {
	case *Started:
		state.replyTo.Tell(EchoResponse{})
	}
}

func NewEchoOnStartActor(replyTo *process.ID) func() Actor {
	return func() Actor {
		return &EchoOnStartActor{replyTo: replyTo}
	}
}

func TestActorCanReplyOnStarting(t *testing.T) {
	future := process.NewFuture(testTimeout)
	actor := Spawn(FromProducer(NewEchoOnStartActor(future.PID())))
	defer StopActor(actor)
	if _, err := future.Result(); err != nil {
		assert.Fail(t, "timed out")
		return
	}
}

type EchoOnStoppingActor struct{ replyTo *process.ID }

func (state *EchoOnStoppingActor) Receive(context Context) {
	switch context.Message().(type) {
	case *Stopping:
		state.replyTo.Tell(EchoResponse{})
	}
}

func NewEchoOnStoppingActor(replyTo *process.ID) func() Actor {
	return func() Actor {
		return &EchoOnStoppingActor{replyTo: replyTo}
	}
}

func TestActorCanReplyOnStopping(t *testing.T) {
	future := process.NewFuture(testTimeout)
	actor := Spawn(FromProducer(NewEchoOnStoppingActor(future.PID())))
	StopActor(actor)
	if _, err := future.Result(); err != nil {
		assert.Fail(t, "timed out")
		return
	}
}
