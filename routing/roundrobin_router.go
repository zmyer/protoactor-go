package routing

import (
	"sync/atomic"

	"github.com/AsynkronIT/protoactor-go/process"
)

type RoundRobinGroupRouter struct {
	GroupRouter
}

type RoundRobinPoolRouter struct {
	PoolRouter
}

type RoundRobinState struct {
	index   int32
	routees *process.PIDSet
	values  []process.PID
}

func (state *RoundRobinState) SetRoutees(routees *process.PIDSet) {
	state.routees = routees
	state.values = routees.Values()
}

func (state *RoundRobinState) GetRoutees() *process.PIDSet {
	return state.routees
}

func (state *RoundRobinState) RouteMessage(message interface{}, sender *process.PID) {
	pid := roundRobinRoutee(&state.index, state.values)
	pid.Request(message, sender)
}

func NewRoundRobinGroup(routees ...*process.PID) GroupRouterConfig {
	r := &RoundRobinGroupRouter{}
	r.Routees = process.NewPIDSet(routees...)
	return r
}

func NewRoundRobinPool(poolSize int) PoolRouterConfig {
	r := &RoundRobinPoolRouter{}
	r.PoolSize = poolSize
	return r
}

func (config *RoundRobinPoolRouter) CreateRouterState() RouterState {
	return &RoundRobinState{}
}

func (config *RoundRobinGroupRouter) CreateRouterState() RouterState {
	return &RoundRobinState{}
}

func roundRobinRoutee(index *int32, routees []process.PID) process.PID {
	i := int(atomic.AddInt32(index, 1))
	mod := len(routees)
	routee := routees[i%mod]
	return routee
}
