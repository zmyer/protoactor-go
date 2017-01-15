package routing

import (
	"log"

	"github.com/AsynkronIT/protoactor-go/process"
	"github.com/serialx/hashring"
)

type Hasher interface {
	Hash() string
}

type ConsistentHashGroupRouter struct {
	GroupRouter
}

type ConsistentHashPoolRouter struct {
	PoolRouter
}

type ConsistentHashRouterState struct {
	hashring  *hashring.HashRing
	routeeMap map[string]*process.ID
}

func (state *ConsistentHashRouterState) SetRoutees(routees *process.PIDSet) {
	//lookup from node name to PID
	state.routeeMap = make(map[string]*process.ID)
	nodes := make([]string, routees.Len())
	routees.ForEach(func(i int, pid process.ID) {
		nodeName := pid.Address + "@" + pid.Id
		nodes[i] = nodeName
		state.routeeMap[nodeName] = &pid
	})
	//initialize hashring for mapping message keys to node names
	state.hashring = hashring.New(nodes)
}

func (state *ConsistentHashRouterState) GetRoutees() *process.PIDSet {
	var routees process.PIDSet
	for _, v := range state.routeeMap {
		routees.Add(v)
	}
	return &routees
}

func (state *ConsistentHashRouterState) RouteMessage(message interface{}, sender *process.ID) {
	switch msg := message.(type) {
	case Hasher:
		key := msg.Hash()

		node, ok := state.hashring.GetNode(key)
		if !ok {
			log.Printf("[ROUTING] Consistent has router failed to derminate routee: %v", key)
			return
		}
		if routee, ok := state.routeeMap[node]; ok {
			routee.Request(msg, sender)
		} else {
			log.Println("[ROUTING] Consisten router failed to resolve node", node)
		}
	default:
		log.Println("[ROUTING] Message must implement routing.Hasher", msg)
	}
}

func (state *ConsistentHashRouterState) InvokeRouterManagementMessage(msg ManagementMessage, sender *process.ID) {

}

func NewConsistentHashPool(poolSize int) PoolRouterConfig {
	r := &ConsistentHashPoolRouter{}
	r.PoolSize = poolSize
	return r
}

func NewConsistentHashGroup(routees ...*process.ID) GroupRouterConfig {
	r := &ConsistentHashGroupRouter{}
	r.Routees = process.NewPIDSet(routees...)
	return r
}

func (config *ConsistentHashPoolRouter) CreateRouterState() RouterState {
	return &ConsistentHashRouterState{}
}

func (config *ConsistentHashGroupRouter) CreateRouterState() RouterState {
	return &ConsistentHashRouterState{}
}
