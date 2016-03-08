package mock

import "github.com/drone/drone/bus"
import "github.com/stretchr/testify/mock"

type Bus struct {
	mock.Mock
}

func (_m *Bus) Publish(_a0 *bus.Event) {
	_m.Called(_a0)
}
func (_m *Bus) Subscribe(_a0 chan *bus.Event) {
	_m.Called(_a0)
}
func (_m *Bus) Unsubscribe(_a0 chan *bus.Event) {
	_m.Called(_a0)
}
