package process

//SystemMessage is a special type of messages passed to control the actor lifecycles
type SystemMessage interface {
	SystemMessage()
}
