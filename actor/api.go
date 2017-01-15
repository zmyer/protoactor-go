package actor

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/process"
)

func StopActor(pid *process.PID) {
	sendSystemMessage(pid, stopMessage)
}

func StopActorFuture(pid *process.PID) *process.Future {
	future := process.NewFuture(10 * time.Second)

	sendSystemMessage(pid, &Watch{Watcher: future.PID()})
	StopActor(pid)

	return future
}

func sendSystemMessage(pid *process.PID, msg process.SystemMessage) {
	s, _ := process.Registry.Get(pid)
	s.SendSystemMessage(pid, msg)
}
