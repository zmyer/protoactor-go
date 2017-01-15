package streams

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/process"
)

type UntypedStream struct {
	c   chan interface{}
	pid *process.ID
}

func (s *UntypedStream) C() <-chan interface{} {
	return s.c
}

func (s *UntypedStream) PID() *process.ID {
	return s.pid
}

func (s *UntypedStream) Close() {
	actor.StopActor(s.pid)
	close(s.c)
}

func NewUntypedStream() *UntypedStream {
	c := make(chan interface{})
	props := actor.FromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		case actor.AutoReceiveMessage:
		case actor.SystemMessage: //ignore terminate
		default:
			c <- msg
		}
	})
	pid := actor.Spawn(props)

	stream := &UntypedStream{
		c:   c,
		pid: pid,
	}
	return stream
}
