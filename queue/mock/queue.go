package mock

import "github.com/drone/drone/queue"
import "github.com/stretchr/testify/mock"

type Queue struct {
	mock.Mock
}

func (_m *Queue) Publish(_a0 *queue.Work) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*queue.Work) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *Queue) Remove(_a0 *queue.Work) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*queue.Work) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *Queue) Pull() *queue.Work {
	ret := _m.Called()

	var r0 *queue.Work
	if rf, ok := ret.Get(0).(func() *queue.Work); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*queue.Work)
		}
	}

	return r0
}
func (_m *Queue) PullClose(_a0 queue.CloseNotifier) *queue.Work {
	ret := _m.Called(_a0)

	var r0 *queue.Work
	if rf, ok := ret.Get(0).(func(queue.CloseNotifier) *queue.Work); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*queue.Work)
		}
	}

	return r0
}
