package process

import (
	"log"
	"strings"
	"time"
)

//Tell a message to a given PID
func (pid *ID) Tell(message interface{}) {
	ref, _ := Registry.Get(pid)
	ref.SendUserMessage(pid, message, nil)
}

//Ask a message to a given PID
func (pid *ID) Request(message interface{}, respondTo *ID) {
	ref, _ := Registry.Get(pid)
	ref.SendUserMessage(pid, message, respondTo)
}

//RequestFuture sends a message to a given PID and returns a Future
func (pid *ID) RequestFuture(message interface{}, timeout time.Duration) *Future {
	ref, ok := Registry.Get(pid)
	if !ok {
		log.Printf("[ACTOR] RequestFuture for missing local PID '%v'", pid.String())
	}

	future := NewFuture(timeout)
	ref.SendUserMessage(pid, message, future.PID())
	return future
}

func pidFromKey(key string, p *ID) {
	i := strings.IndexByte(key, '#')
	if i == -1 {
		p.Address = Registry.Address
		p.Id = key
	} else {
		p.Address = key[:i]
		p.Id = key[i+1:]
	}
}

func (pid *ID) key() string {
	if pid.Address == Registry.Address {
		return pid.Id
	}
	return pid.Address + "#" + pid.Id
}

func (pid *ID) Empty() bool {
	return pid.Address == "" && pid.Id == ""
}

func (pid *ID) String() string {
	return pid.Address + "/" + pid.Id
}

//NewPID returns a new instance of the PID struct
func NewPID(address, id string) *ID {
	return &ID{
		Address: address,
		Id:      id,
	}
}

//NewLocalPID returns a new instance of the PID struct with the address preset
func NewLocalPID(id string) *ID {
	return &ID{
		Address: Registry.Address,
		Id:      id,
	}
}
