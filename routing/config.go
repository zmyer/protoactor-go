package routing

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/process"
)

type RouterConfig interface {
	OnStarted(context actor.Context, props actor.Props, router RouterState)
	CreateRouterState() RouterState
}

type GroupRouterConfig interface {
	RouterConfig
}

type PoolRouterConfig interface {
	RouterConfig
}

type GroupRouter struct {
	RouterConfig
	Routees *process.PIDSet
}

type PoolRouter struct {
	RouterConfig
	PoolSize int
}

func (config *GroupRouter) OnStarted(context actor.Context, props actor.Props, router RouterState) {
	config.Routees.ForEach(func(i int, pid process.PID) {
		context.Watch(&pid)
	})
	router.SetRoutees(config.Routees)
}

func (config *PoolRouter) OnStarted(context actor.Context, props actor.Props, router RouterState) {
	var routees process.PIDSet
	for i := 0; i < config.PoolSize; i++ {
		routees.Add(context.Spawn(props))
	}
	router.SetRoutees(&routees)
}

func spawner(config RouterConfig) actor.Spawner {
	return func(id string, props actor.Props, parent *process.PID) *process.PID {
		return spawn(id, config, props, parent)
	}
}

func FromProps(props actor.Props, config RouterConfig) actor.Props {
	return props.WithSpawn(spawner(config))
}

func FromGroupRouter(config GroupRouterConfig) actor.Props {
	return actor.Props{}.WithSpawn(spawner(config))
}
