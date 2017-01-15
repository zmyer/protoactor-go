package process

type deadLetterProcess struct{}

var (
	deadLetter Process = &deadLetterProcess{}
)

type DeadLetter struct {
	PID     *ID
	Message interface{}
	Sender  *ID
}

func (*deadLetterProcess) SendUserMessage(pid *ID, message interface{}, sender *ID) {
	EventStream.Publish(&DeadLetter{
		PID:     pid,
		Message: message,
		Sender:  sender,
	})
}

func (*deadLetterProcess) SendSystemMessage(pid *ID, message SystemMessage) {
	EventStream.Publish(&DeadLetter{
		PID:     pid,
		Message: message,
	})
}
