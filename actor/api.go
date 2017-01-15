package actor

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/process"
)

func stopActor(pid *process.PID) {
	sendSystemMessage(pid, stopMessage)
}

func stopActorFuture(pid *process.PID) *process.Future {
	future := process.NewFuture(10 * time.Second)

	sendSystemMessage(pid, &Watch{Watcher: future.PID()})
	stopActor(pid)

	return future
}

func sendSystemMessage(pid *process.PID, msg process.SystemMessage) {
	s, _ := process.ProcessRegistry.Get(pid)
	s.SendSystemMessage(pid, msg)
}
