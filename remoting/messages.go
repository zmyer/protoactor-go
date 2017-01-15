package remoting

import "github.com/AsynkronIT/protoactor-go/process"

type EndpointTerminated struct {
	Address string
}

type remoteWatch struct {
	Watcher *process.PID
	Watchee *process.PID
}

type remoteUnwatch struct {
	Watcher *process.PID
	Watchee *process.PID
}

type remoteTerminate struct {
	Watcher *process.PID
	Watchee *process.PID
}
