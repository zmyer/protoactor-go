package actor

type rootContext struct {
	meta ActorMeta
}

func NewRootContext() MetaContext {
	return &rootContext{
		meta: make(map[string]string),
	}
}

type ActorMeta map[string]string

func (a ActorMeta) Add(key, value string) {
	a[key] = value
}

func (a ActorMeta) Get(key string) string {
	return a[key]
}

type MetaContext interface {
	Tell(message interface{}, target *PID)
	Meta() ActorMeta
}

func (rc *rootContext) Tell(message interface{}, target *PID) {
	target.Tell(message)
}
func (rc *rootContext) Meta() ActorMeta {
	return rc.meta
}
