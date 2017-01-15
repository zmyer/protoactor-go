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

func spawnMockProcess(name string) (*process.ID, *mockProcess) {
	p := &mockProcess{}
	pid, ok := process.Registry.Add(p, name)
	if !ok {
		panic(fmt.Errorf("did not spawn named process '%s'", name))
	}

	return pid, p
}

func removeMockProcess(pid *process.ID) {
	process.Registry.Remove(pid)
}

type mockProcess struct {
	mock.Mock
}

func (m *mockProcess) SendUserMessage(pid *process.ID, message interface{}, sender *process.ID) {
	m.Called(pid, message, sender)
}
func (m *mockProcess) SendSystemMessage(pid *process.ID, message process.SystemMessage) {
	m.Called(pid, message)
}
func (m *mockProcess) Stop(pid *process.ID) {
	m.Called(pid)
}

type mockContext struct {
	mock.Mock
}

func (m *mockContext) Watch(pid *process.ID) {
	m.Called(pid)
}

func (m *mockContext) Unwatch(pid *process.ID) {
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

func (m *mockContext) Sender() *process.ID {
	args := m.Called()
	return args.Get(0).(*process.ID)
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

func (m *mockContext) Self() *process.ID {
	args := m.Called()
	return args.Get(0).(*process.ID)
}

func (m *mockContext) Parent() *process.ID {
	args := m.Called()
	return args.Get(0).(*process.ID)
}

func (m *mockContext) Spawn(p Props) *process.ID {
	args := m.Called(p)
	return args.Get(0).(*process.ID)
}

func (m *mockContext) SpawnNamed(p Props, name string) *process.ID {
	args := m.Called(p, name)
	return args.Get(0).(*process.ID)
}

func (m *mockContext) Children() []*process.ID {
	args := m.Called()
	return args.Get(0).([]*process.ID)
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
