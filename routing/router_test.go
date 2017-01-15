package routing

import (
	"fmt"
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/process"
	"github.com/stretchr/testify/mock"
)

var _ fmt.Formatter
var _ time.Time

func TestRouterSendsUserMessageToChild(t *testing.T) {
	child, p := spawnMockProcess("child")
	defer removeMockProcess(child)

	p.On("SendUserMessage", mock.Anything, "hello", mock.Anything)
	p.On("SendSystemMessage", mock.Anything, mock.Anything, mock.Anything)

	s1 := process.NewPIDSet(child)

	rs := new(testRouterState)
	rs.On("SetRoutees", s1)
	rs.On("RouteMessage", "hello", mock.Anything)

	grc := newGroupRouterConfig(child)
	grc.On("CreateRouterState").Return(rs)

	routerPID := SpawnGroup(grc)
	routerPID.Tell("hello")

	mock.AssertExpectationsForObjects(t, p, rs)
}

type testGroupRouter struct {
	GroupRouter
	mock.Mock
}

func newGroupRouterConfig(routees ...*process.ID) *testGroupRouter {
	r := new(testGroupRouter)
	r.Routees = process.NewPIDSet(routees...)
	return r
}

func (m *testGroupRouter) CreateRouterState() RouterState {
	fmt.Println("Doing it")
	args := m.Called()
	return args.Get(0).(*testRouterState)
}

type testRouterState struct {
	mock.Mock
	routees *process.PIDSet
}

func (m *testRouterState) SetRoutees(routees *process.PIDSet) {
	fmt.Println("SetRoutees")
	m.Called(routees)
	m.routees = routees
}

func (m *testRouterState) RouteMessage(message interface{}, sender *process.ID) {
	m.Called(message, sender)
	m.routees.ForEach(func(i int, pid process.ID) {
		pid.Request(message, sender)
	})
}

func (m *testRouterState) GetRoutees() *process.PIDSet {
	args := m.Called()
	return args.Get(0).(*process.PIDSet)
}
