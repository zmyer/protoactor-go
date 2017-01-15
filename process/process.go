package process

// Process is an interface that defines the base contract for interaction of actors
type Process interface {
	SendUserMessage(pid *ID, message interface{}, sender *ID)
	SendSystemMessage(pid *ID, message SystemMessage)
}
