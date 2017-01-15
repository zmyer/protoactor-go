package routing

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/process"
)

type routerActor struct {
	props  actor.Props
	config RouterConfig
	state  RouterState
}

func (a *routerActor) Receive(context actor.Context) {
	switch m := context.Message().(type) {
	case *actor.Started:
		a.config.OnStarted(context, a.props, a.state)

	case *AddRoutee:
		r := a.state.GetRoutees()
		if r.Contains(m.PID) {
			return
		}
		context.Watch(m.PID)
		r.Add(m.PID)
		a.state.SetRoutees(r)

	case *RemoveRoutee:
		r := a.state.GetRoutees()
		if !r.Contains(m.PID) {
			return
		}

		context.Unwatch(m.PID)
		r.Remove(m.PID)
		a.state.SetRoutees(r)

	case *BroadcastMessage:
		msg := m.Message
		sender := context.Sender()
		a.state.GetRoutees().ForEach(func(i int, pid process.ID) {
			pid.Request(msg, sender)
		})

	case *GetRoutees:
		r := a.state.GetRoutees()
		routees := make([]*process.ID, r.Len())
		r.ForEach(func(i int, pid process.ID) {
			routees[i] = &pid
		})

		context.Sender().Tell(&Routees{routees})
	}
}
