package actor

import "github.com/AsynkronIT/protoactor-go/process"

type Directive int

const (
	ResumeDirective Directive = iota
	RestartDirective
	StopDirective
	EscalateDirective
)

type Decider func(child *process.ID, cause interface{}) Directive

//TODO: as we dont allow remote children or remote SupervisionStrategy
//Instead of letting the parent keep track of child restart stats.
//this info could actually go into each actor, sending it back to the parent as part of the Failure message
type SupervisorStrategy interface {
	HandleFailure(supervisor Supervisor, child *process.ID, cause interface{})
}

type OneForOneStrategy struct {
	maxNrOfRetries              int
	withinTimeRangeMilliseconds int
	decider                     Decider
}

type Supervisor interface {
	Children() []*process.ID
	EscalateFailure(who *process.ID, reason interface{})
}

func (strategy *OneForOneStrategy) HandleFailure(supervisor Supervisor, child *process.ID, reason interface{}) {
	directive := strategy.decider(child, reason)

	switch directive {
	case ResumeDirective:
		//resume the failing child
		sendSystemMessage(child, resumeMailboxMessage)
	case RestartDirective:
		//restart the failing child
		sendSystemMessage(child, restartMessage)
	case StopDirective:
		//stop the failing child
		StopActor(child)
	case EscalateDirective:
		//send failure to parent
		//supervisor mailbox
		supervisor.EscalateFailure(child, reason)
	}
}

func NewOneForOneStrategy(maxNrOfRetries int, withinTimeRangeMilliseconds int, decider Decider) SupervisorStrategy {
	return &OneForOneStrategy{
		maxNrOfRetries:              maxNrOfRetries,
		withinTimeRangeMilliseconds: withinTimeRangeMilliseconds,
		decider:                     decider,
	}
}

func DefaultDecider(child *process.ID, reason interface{}) Directive {
	return RestartDirective
}

var defaultSupervisionStrategy = NewOneForOneStrategy(10, 3000, DefaultDecider)

func DefaultSupervisionStrategy() SupervisorStrategy {
	return defaultSupervisionStrategy
}
