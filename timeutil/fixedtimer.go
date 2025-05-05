package timeutil

import (
	"sync"
	"time"
)

// FixedTimer implements a timer that triggers a task function at fixed regular intervals
type FixedTimer struct {
	mu       sync.Mutex
	running  bool
	interval time.Duration
	tick     TickFn
	ticker   *time.Ticker
	stop     chan struct{}
	wg       sync.WaitGroup
}

// Start starts the timer and executes the provided function at each interval
func (t *FixedTimer) Start() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.running {
		return
	}

	t.running = true
	t.ticker = time.NewTicker(t.interval)
	t.wg.Add(1)

	go func() {
		defer t.wg.Done()
		for {
			select {
			case <-t.ticker.C:
				t.tick()
			case <-t.stop:
				return
			}
		}
	}()
}

// Stop stops the timer and cleans up resources
func (t *FixedTimer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.running {
		return
	}

	t.ticker.Stop()
	close(t.stop)
	t.wg.Wait()
	t.running = false

	// prepare to reuse stop chan again
	t.stop = make(chan struct{})
}

// NewFixedTimer returns new Timer value
func NewFixedTimer(dur time.Duration, fn TickFn) Timer {
	return &FixedTimer{
		interval: dur,
		tick:     fn,
		stop:     make(chan struct{}),
	}
}
