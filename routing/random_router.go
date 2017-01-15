package routing

import (
	"math/rand"

	"github.com/AsynkronIT/protoactor-go/process"
)

type RandomGroupRouter struct {
	GroupRouter
}

type RandomPoolRouter struct {
	PoolRouter
}

type RandomRouterState struct {
	routees *process.PIDSet
	values  []process.ID
}

func (state *RandomRouterState) SetRoutees(routees *process.PIDSet) {
	state.routees = routees
	state.values = routees.Values()
}

func (state *RandomRouterState) GetRoutees() *process.PIDSet {
	return state.routees
}

func (state *RandomRouterState) RouteMessage(message interface{}, sender *process.ID) {
	l := len(state.values)
	r := rand.Intn(l)
	pid := state.values[r]
	pid.Request(message, sender)
}

func NewRandomPool(poolSize int) PoolRouterConfig {
	r := &RandomPoolRouter{}
	r.PoolSize = poolSize
	return r
}

func NewRandomGroup(routees ...*process.ID) GroupRouterConfig {
	r := &RandomGroupRouter{}
	r.Routees = process.NewPIDSet(routees...)
	return r
}

func (config *RandomPoolRouter) CreateRouterState() RouterState {
	return &RandomRouterState{}
}

func (config *RandomGroupRouter) CreateRouterState() RouterState {
	return &RandomRouterState{}
}
