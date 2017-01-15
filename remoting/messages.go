package remoting

import "github.com/AsynkronIT/protoactor-go/process"

type EndpointTerminated struct {
	Address string
}

type remoteWatch struct {
	Watcher *process.ID
	Watchee *process.ID
}

type remoteUnwatch struct {
	Watcher *process.ID
	Watchee *process.ID
}

type remoteTerminate struct {
	Watcher *process.ID
	Watchee *process.ID
}
