package routing

import "github.com/AsynkronIT/protoactor-go/process"

type BroadcastGroupRouter struct {
	GroupRouter
}

type BroadcastPoolRouter struct {
	PoolRouter
}

type BroadcastRouterState struct {
	routees *process.PIDSet
}

func (state *BroadcastRouterState) SetRoutees(routees *process.PIDSet) {
	state.routees = routees
}

func (state *BroadcastRouterState) GetRoutees() *process.PIDSet {
	return state.routees
}

func (state *BroadcastRouterState) RouteMessage(message interface{}, sender *process.PID) {
	state.routees.ForEach(func(i int, pid process.PID) {
		pid.Request(message, sender)
	})
}

func NewBroadcastPool(poolSize int) PoolRouterConfig {
	r := &BroadcastPoolRouter{}
	r.PoolSize = poolSize
	return r
}

func NewBroadcastGroup(routees ...*process.PID) GroupRouterConfig {
	r := &BroadcastGroupRouter{}
	r.Routees = process.NewPIDSet(routees...)
	return r
}

func (config *BroadcastPoolRouter) CreateRouterState() RouterState {
	return &BroadcastRouterState{}
}

func (config *BroadcastGroupRouter) CreateRouterState() RouterState {
	return &BroadcastRouterState{}
}
