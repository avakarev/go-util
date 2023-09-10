package timeutil

import (
	"sync"
	"time"
)

// Clock defines interface to the time functions
type Clock interface {
	Now() time.Time
	Since(t time.Time) time.Duration
	Until(t time.Time) time.Duration
}

// clock adopts Clock to stdlib time
type clock struct{}

// Now adopts time.Now
func (c *clock) Now() time.Time {
	return time.Now()
}

// Since adopts time.Since
func (c *clock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// Until adopts time.Until
func (c *clock) Until(t time.Time) time.Duration {
	return time.Until(t)
}

// NewClock returns an new clock value
func NewClock() Clock {
	return &clock{}
}

// Mock implements clock mock that can adjust time on demand
type Mock struct {
	mu  sync.Mutex
	now time.Time
}

// Now returns the current wall time on the mock clock
func (m *Mock) Now() time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.now
}

// Since returns the time elapsed since t
func (m *Mock) Since(t time.Time) time.Duration {
	return m.Now().Sub(t)
}

// Until returns the duration until t
func (m *Mock) Until(t time.Time) time.Duration {
	return t.Sub(m.Now())
}

// Add adjusts the current time of the mock clock by given duration
func (m *Mock) Add(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.now = m.now.Add(d)
}

// Set sets the current time of the mock clock to a given value
func (m *Mock) Set(t time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.now = t
}

// NewMock returns new clock mock value
func NewMock() *Mock {
	return &Mock{now: time.Unix(0, 0)}
}
