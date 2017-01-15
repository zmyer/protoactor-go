package actor

import "github.com/AsynkronIT/protoactor-go/process"

type localProcess struct {
	mailbox Mailbox
}

func newLocalProcess(mailbox Mailbox) *localProcess {
	return &localProcess{
		mailbox: mailbox,
	}
}

func (ref *localProcess) SendUserMessage(pid *process.PID, message interface{}, sender *process.PID) {
	if sender != nil {
		ref.mailbox.PostUserMessage(&messageSender{Message: message, Sender: sender})
	} else {
		ref.mailbox.PostUserMessage(message)
	}
}

func (ref *localProcess) SendSystemMessage(pid *process.PID, message process.SystemMessage) {
	ref.mailbox.PostSystemMessage(message)
}

func (ref *localProcess) Stop(pid *process.PID) {
	ref.SendSystemMessage(pid, stopMessage)
}
