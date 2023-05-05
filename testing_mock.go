package expect

import (
	"sync/atomic"
)

// TMock is a mock implementation of the T
// interface.
type TMock struct {
	HelperStub   func()
	HelperCalled int32
	ErrorfStub   func(format string, args ...any)
	ErrorfCalled int32
}

var _ T = &TMock{}

// Helper is a stub for the T.Helper
// method that records the number of times it has been called.
func (m *TMock) Helper() {
	atomic.AddInt32(&m.HelperCalled, 1)
	m.HelperStub()
}

// Errorf is a stub for the T.Errorf
// method that records the number of times it has been called.
func (m *TMock) Errorf(format string, args ...any) {
	atomic.AddInt32(&m.ErrorfCalled, 1)
	m.ErrorfStub(format, args...)
}
