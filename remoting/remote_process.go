package remoting

import (
	"log"
	"reflect"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/process"
	"github.com/gogo/protobuf/proto"
)

type remoteProcess struct {
	pid *process.PID
}

func newRemoteProcess(pid *process.PID) process.Process {
	return &remoteProcess{
		pid: pid,
	}
}

func (ref *remoteProcess) SendUserMessage(pid *process.PID, message interface{}, sender *process.PID) {
	sendRemoteMessage(pid, message, sender)
}

func sendRemoteMessage(pid *process.PID, message interface{}, sender *process.PID) {
	switch msg := message.(type) {
	case proto.Message:
		envelope, _ := serialize(msg, pid, sender)
		endpointManagerPID.Tell(envelope)
	default:
		log.Printf("[REMOTING] failed, trying to send non Proto %s message to %v", reflect.TypeOf(msg), pid)
	}
}

func (ref *remoteProcess) SendSystemMessage(pid *process.PID, message process.SystemMessage) {

	//intercept any Watch messages and direct them to the endpoint manager
	switch msg := message.(type) {
	case *actor.Watch:
		rw := &remoteWatch{
			Watcher: msg.Watcher,
			Watchee: pid,
		}
		endpointManagerPID.Tell(rw)
	case *actor.Unwatch:
		ruw := &remoteUnwatch{
			Watcher: msg.Watcher,
			Watchee: pid,
		}
		endpointManagerPID.Tell(ruw)
	default:
		sendRemoteMessage(pid, message, nil)
	}
}

func (ref *remoteProcess) Stop(pid *process.PID) {
	ref.SendSystemMessage(pid, &actor.Stop{})
}
