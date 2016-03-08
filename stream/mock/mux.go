package mock

import "github.com/stretchr/testify/mock"

import "io"

type Mux struct {
	mock.Mock
}

func (_m *Mux) Create(key string) (io.ReadCloser, io.WriteCloser, error) {
	ret := _m.Called(key)

	var r0 io.ReadCloser
	if rf, ok := ret.Get(0).(func(string) io.ReadCloser); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(io.ReadCloser)
	}

	var r1 io.WriteCloser
	if rf, ok := ret.Get(1).(func(string) io.WriteCloser); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Get(1).(io.WriteCloser)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(key)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
func (_m *Mux) Open(key string) (io.ReadCloser, io.WriteCloser, error) {
	ret := _m.Called(key)

	var r0 io.ReadCloser
	if rf, ok := ret.Get(0).(func(string) io.ReadCloser); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(io.ReadCloser)
	}

	var r1 io.WriteCloser
	if rf, ok := ret.Get(1).(func(string) io.WriteCloser); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Get(1).(io.WriteCloser)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(key)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
func (_m *Mux) Remove(key string) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *Mux) Exists(key string) bool {
	ret := _m.Called(key)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
