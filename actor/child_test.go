package actor

import (
	"testing"

	"time"

	"github.com/AsynkronIT/protoactor-go/process"
	"github.com/stretchr/testify/assert"
)

type CreateChildMessage struct{}
type GetChildCountRequest struct{}
type GetChildCountResponse struct{ ChildCount int }
type CreateChildActor struct{}

func (*CreateChildActor) Receive(context Context) {
	switch context.Message().(type) {
	case CreateChildMessage:
		context.Spawn(FromProducer(NewBlackHoleActor))
	case GetChildCountRequest:
		reply := GetChildCountResponse{ChildCount: len(context.Children())}
		context.Respond(reply)
	}
}

func NewCreateChildActor() Actor {
	return &CreateChildActor{}
}

func TestActorCanCreateChildren(t *testing.T) {
	actor := Spawn(FromProducer(NewCreateChildActor))
	defer StopActor(actor)
	expected := 10
	for i := 0; i < expected; i++ {
		actor.Tell(CreateChildMessage{})
	}
	response, err := actor.RequestFuture(GetChildCountRequest{}, testTimeout).Result()
	if err != nil {
		assert.Fail(t, "timed out")
		return
	}
	assert.Equal(t, expected, response.(GetChildCountResponse).ChildCount)
}

type CreateChildThenStopActor struct {
	replyTo *process.ID
}

type GetChildCountMessage2 struct {
	ReplyDirectly  *process.ID
	ReplyAfterStop *process.ID
}

func (state *CreateChildThenStopActor) Receive(context Context) {
	switch msg := context.Message().(type) {
	case CreateChildMessage:
		context.Spawn(FromProducer(NewBlackHoleActor))
	case GetChildCountMessage2:
		msg.ReplyDirectly.Tell(true)
		state.replyTo = msg.ReplyAfterStop
	case *Stopped:
		reply := GetChildCountResponse{ChildCount: len(context.Children())}
		state.replyTo.Tell(reply)
	}
}

func NewCreateChildThenStopActor() Actor {
	return &CreateChildThenStopActor{}
}

func TestActorCanStopChildren(t *testing.T) {

	actor := Spawn(FromProducer(NewCreateChildThenStopActor))
	count := 10
	for i := 0; i < count; i++ {
		actor.Tell(CreateChildMessage{})
	}

	future := process.NewFuture(testTimeout)
	future2 := process.NewFuture(testTimeout)
	actor.Tell(GetChildCountMessage2{ReplyDirectly: future.PID(), ReplyAfterStop: future2.PID()})

	//wait for the actor to reply to the first responsePID
	err := future.Wait()
	if err != nil {
		assert.Fail(t, "timed out")
		return
	}

	//then send a stop command
	StopActor(actor)

	//wait for the actor to stop and get the result from the stopped handler
	response, err := future2.Result()
	if err != nil {
		assert.Fail(t, "timed out")
		return
	}
	//we should have 0 children when the actor is stopped
	assert.Equal(t, 0, response.(GetChildCountResponse).ChildCount)
}

func TestActorReceivesTerminatedFromWatched(t *testing.T) {
	child := Spawn(FromInstance(nullReceive))
	future := process.NewFuture(testTimeout)
	var r Receive = func(c Context) {
		switch msg := c.Message().(type) {
		case *Started:
			c.Watch(child)

		case *Terminated:
			ac := c.(*actorCell)
			if msg.Who.Equal(child) && ac.watching.Empty() {
				future.PID().Tell(true)
			}
		}
	}

	Spawn(FromInstance(r))
	StopActor(child)

	_, err := future.Result()
	assert.NoError(t, err, "timed out")
}

func TestFutureDoesTimeout(t *testing.T) {
	pid := Spawn(FromInstance(nullReceive))
	_, err := pid.RequestFuture("", time.Millisecond).Result()
	assert.EqualError(t, err, process.ErrTimeout.Error())
}

func TestFutureDoesNotTimeout(t *testing.T) {
	var r Receive = func(c Context) {
		if _, ok := c.Message().(string); !ok {
			return
		}

		time.Sleep(50 * time.Millisecond)
		c.Respond("foo")
	}
	pid := Spawn(FromInstance(r))
	reply, err := pid.RequestFuture("", 2*time.Second).Result()
	assert.NoError(t, err)
	assert.Equal(t, "foo", reply)
}
