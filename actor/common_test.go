package actor

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/AsynkronIT/protoactor-go/process"
	"github.com/stretchr/testify/mock"
)

var nullReceive Receive = func(Context) {}

func init() {
	// discard all logging in tests
	log.SetOutput(ioutil.Discard)
}

func spawnMockProcess(name string) (*process.PID, *mockProcess) {
	p := &mockProcess{}
	pid, ok := process.Registry.Add(p, name)
	if !ok {
		panic(fmt.Errorf("did not spawn named process '%s'", name))
	}

	return pid, p
}

func removeMockProcess(pid *process.PID) {
	process.Registry.Remove(pid)
}

type mockProcess struct {
	mock.Mock
}

func (m *mockProcess) SendUserMessage(pid *process.PID, message interface{}, sender *process.PID) {
	m.Called(pid, message, sender)
}
func (m *mockProcess) SendSystemMessage(pid *process.PID, message process.SystemMessage) {
	m.Called(pid, message)
}
func (m *mockProcess) Stop(pid *process.PID) {
	m.Called(pid)
}

type mockContext struct {
	mock.Mock
}

func (m *mockContext) Watch(pid *process.PID) {
	m.Called(pid)
}

func (m *mockContext) Unwatch(pid *process.PID) {
	m.Called(pid)
}

func (m *mockContext) Message() interface{} {
	args := m.Called()
	return args.Get(0)
}

func (m *mockContext) SetReceiveTimeout(d time.Duration) {
	m.Called(d)
}
func (m *mockContext) ReceiveTimeout() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func (m *mockContext) Sender() *process.PID {
	args := m.Called()
	return args.Get(0).(*process.PID)
}

func (m *mockContext) Become(r Receive) {
	m.Called(r)
}

func (m *mockContext) BecomeStacked(r Receive) {
	m.Called(r)
}

func (m *mockContext) UnbecomeStacked() {
	m.Called()
}

func (m *mockContext) Self() *process.PID {
	args := m.Called()
	return args.Get(0).(*process.PID)
}

func (m *mockContext) Parent() *process.PID {
	args := m.Called()
	return args.Get(0).(*process.PID)
}

func (m *mockContext) Spawn(p Props) *process.PID {
	args := m.Called(p)
	return args.Get(0).(*process.PID)
}

func (m *mockContext) SpawnNamed(p Props, name string) *process.PID {
	args := m.Called(p, name)
	return args.Get(0).(*process.PID)
}

func (m *mockContext) Children() []*process.PID {
	args := m.Called()
	return args.Get(0).([]*process.PID)
}

func (m *mockContext) Next() {
	m.Called()
}

func (m *mockContext) Receive(i interface{}) {
	m.Called(i)
}

func (m *mockContext) Stash() {
	m.Called()
}

func (m *mockContext) Respond(response interface{}) {
	m.Called(response)
}

func (m *mockContext) Actor() Actor {
	args := m.Called()
	return args.Get(0).(Actor)
}
