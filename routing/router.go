package routing

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/process"
)

type routerProcess struct {
	router *process.PID
	state  RouterState
}

func (ref *routerProcess) SendUserMessage(pid *process.PID, message interface{}, sender *process.PID) {
	if _, ok := message.(ManagementMessage); ok {
		r, _ := process.Registry.Get(ref.router)
		r.SendUserMessage(pid, message, sender)
	} else {
		ref.state.RouteMessage(message, sender)
	}
}

func (ref *routerProcess) SendSystemMessage(pid *process.PID, message process.SystemMessage) {
	r, _ := process.Registry.Get(ref.router)
	r.SendSystemMessage(pid, message)
}

func (ref *routerProcess) Stop(pid *process.PID) {
	ref.SendSystemMessage(pid, &actor.Stop{})
}

type RouterState interface {
	RouteMessage(message interface{}, sender *process.PID)
	SetRoutees(routees *process.PIDSet)
	GetRoutees() *process.PIDSet
}
